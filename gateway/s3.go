// SiYuan Gateway - automatic S3 (MinIO) sync provisioning
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// httpPost 向指定 URL 发送 JSON POST 并返回响应体
func httpPost(url string, body []byte) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// s3Config 内置 MinIO/S3 存储的连接信息，用于自动为每个用户开通 S3 同步
type s3Config struct {
	kernelEndpoint string // 传给内核的端点，带协议，如 http://minio:9000
	hostPort       string // 供 minio-go 使用的 host:port（去掉协议）
	useSSL         bool
	accessKey      string
	secretKey      string
	region         string
	keySecret      string // 用于按用户名确定性派生仓库加密口令
}

// loadS3Config 从环境变量读取内置 S3 配置；未设置 GATEWAY_S3_ENDPOINT 时返回 nil（功能关闭）
func loadS3Config() *s3Config {
	ep := strings.TrimSpace(os.Getenv("GATEWAY_S3_ENDPOINT"))
	if "" == ep {
		return nil
	}

	kernelEndpoint := ep
	if !strings.Contains(ep, "://") {
		kernelEndpoint = "http://" + ep
	}
	host := kernelEndpoint
	useSSL := strings.HasPrefix(kernelEndpoint, "https://")
	host = strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://")
	host = strings.TrimSuffix(host, "/")

	return &s3Config{
		kernelEndpoint: kernelEndpoint,
		hostPort:       host,
		useSSL:         useSSL,
		accessKey:      os.Getenv("GATEWAY_S3_ACCESS_KEY"),
		secretKey:      os.Getenv("GATEWAY_S3_SECRET_KEY"),
		region:         env("GATEWAY_S3_REGION", "us-east-1"),
		keySecret:      os.Getenv("GATEWAY_S3_KEY_SECRET"),
	}
}

// derivePassphrase 由网关密钥与用户名确定性派生仓库加密口令，保证跨设备/重建时同一用户得到同一密钥
func (s *s3Config) derivePassphrase(username string) string {
	secret := s.keySecret
	if "" == secret {
		secret = "siyuan-gateway-default-secret"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("siyuan-repo-key:" + username))
	return hex.EncodeToString(mac.Sum(nil))
}

// ensureBucket 在 MinIO 中确保用户专属 bucket 存在
func (s *s3Config) ensureBucket(name string) error {
	cli, err := minio.New(s.hostPort, &minio.Options{
		Creds:  credentials.NewStaticV4(s.accessKey, s.secretKey, ""),
		Secure: s.useSSL,
		Region: s.region,
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exists, err := cli.BucketExists(ctx, name)
	if err != nil {
		return err
	}
	if !exists {
		return cli.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: s.region})
	}
	return nil
}

// kernelConf 仅解析自动开通所需的字段
type kernelConf struct {
	Sync struct {
		Provider int  `json:"provider"`
		Enabled  bool `json:"enabled"`
		S3       struct {
			Endpoint string `json:"endpoint"`
		} `json:"s3"`
	} `json:"sync"`
	Repo struct {
		Key string `json:"key"`
	} `json:"repo"`
}

// getKernelConf 读取用户内核的当前配置
func (m *KernelManager) getKernelConf(u *User) (*kernelConf, error) {
	body, err := m.kernelPostRaw(u, "/api/system/getConf", map[string]any{})
	if err != nil {
		return nil, err
	}
	var ret struct {
		Data struct {
			Conf kernelConf `json:"conf"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return &ret.Data.Conf, nil
}

// kernelPost 向用户内核发送一次 API 调用并校验返回码
func (m *KernelManager) kernelPost(u *User, path string, payload map[string]any) error {
	body, err := m.kernelPostRaw(u, path, payload)
	if err != nil {
		return err
	}
	var ret struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = json.Unmarshal(body, &ret); err != nil {
		return err
	}
	if 0 != ret.Code {
		return fmt.Errorf("%s failed: %s", path, ret.Msg)
	}
	return nil
}

func (m *KernelManager) kernelPostRaw(u *User, path string, payload map[string]any) ([]byte, error) {
	data, _ := json.Marshal(payload)
	req, err := httpPost(fmt.Sprintf("http://127.0.0.1:%d%s", u.KernelPort, path), data)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ensureS3Sync 首次为用户开通内置 S3 同步：创建 bucket、初始化加密密钥、写入 S3 配置并开启同步。
// 已配置 S3 的用户会被跳过，以免覆盖用户自定义设置。
func (m *KernelManager) ensureS3Sync(u *User) error {
	if nil == m.s3 {
		return nil // 未启用内置 S3
	}

	conf, err := m.getKernelConf(u)
	if err != nil {
		return fmt.Errorf("read kernel conf: %w", err)
	}
	if 2 == conf.Sync.Provider && "" != conf.Sync.S3.Endpoint && conf.Sync.Enabled {
		return nil // 已完整开通并启用，跳过（避免覆盖用户自定义设置）
	}

	// 1. 创建用户专属 bucket
	if err = m.s3.ensureBucket(u.Name); err != nil {
		return fmt.Errorf("ensure bucket: %w", err)
	}

	// 2. 初始化仓库加密密钥（无密钥时）
	if "" == conf.Repo.Key {
		if err = m.kernelPost(u, "/api/repo/initRepoKeyFromPassphrase", map[string]any{"pass": m.s3.derivePassphrase(u.Name)}); err != nil {
			return fmt.Errorf("init repo key: %w", err)
		}
	}

	// 3. 写入 S3 配置
	if err = m.kernelPost(u, "/api/sync/setSyncProviderS3", map[string]any{"s3": map[string]any{
		"endpoint":       m.s3.kernelEndpoint,
		"accessKey":      m.s3.accessKey,
		"secretKey":      m.s3.secretKey,
		"bucket":         u.Name,
		"region":         m.s3.region,
		"pathStyle":      true,
		"skipTlsVerify":  true,
		"timeout":        30,
		"concurrentReqs": 8,
	}}); err != nil {
		return fmt.Errorf("set s3 provider: %w", err)
	}

	// 4. 选择 S3 为同步提供者
	if err = m.kernelPost(u, "/api/sync/setSyncProvider", map[string]any{"provider": 2}); err != nil {
		return fmt.Errorf("set provider: %w", err)
	}

	// 5. 选择云端同步目录名（S3 不支持显式创建目录，首次同步时自动建立）
	if err = m.kernelPost(u, "/api/sync/setCloudSyncDir", map[string]any{"name": "main"}); err != nil {
		return fmt.Errorf("set cloud sync dir: %w", err)
	}

	// 6. 开启同步
	if err = m.kernelPost(u, "/api/sync/setSyncEnable", map[string]any{"enabled": true}); err != nil {
		return fmt.Errorf("enable sync: %w", err)
	}

	log.Printf("s3 sync provisioned for [%s] (bucket [%s])", u.Name, u.Name)
	return nil
}

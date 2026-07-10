// SiYuan Gateway - multi-user front for SiYuan kernels
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// KernelManager 管理每个用户的内核进程生命周期
type KernelManager struct {
	mu        sync.Mutex
	kernelBin string
	dataDir   string
	lang      string
	procs     map[string]*exec.Cmd   // 用户名 -> 内核进程
	starting  map[string]*sync.Mutex // 用户名 -> 启动锁，防止并发重复启动
	s3        *s3Config              // 内置 S3（MinIO）自动同步配置，nil 表示未启用
}

func newKernelManager(kernelBin, dataDir, lang string) *KernelManager {
	return &KernelManager{
		kernelBin: kernelBin,
		dataDir:   dataDir,
		lang:      lang,
		procs:     map[string]*exec.Cmd{},
		starting:  map[string]*sync.Mutex{},
	}
}

func (m *KernelManager) workspacePath(u *User) string {
	return filepath.Join(m.dataDir, u.Name)
}

func (m *KernelManager) userLock(name string) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.starting[name]
	if !ok {
		lock = &sync.Mutex{}
		m.starting[name] = lock
	}
	return lock
}

// Ensure 保证用户内核在运行；若未运行则启动并等待就绪
func (m *KernelManager) Ensure(u *User) error {
	lock := m.userLock(u.Name)
	lock.Lock()
	defer lock.Unlock()

	m.mu.Lock()
	cmd, running := m.procs[u.Name]
	if running && cmd.ProcessState != nil { // 进程已退出
		delete(m.procs, u.Name)
		running = false
	}
	m.mu.Unlock()

	if running || m.isAlive(u) {
		return nil
	}
	return m.start(u)
}

func (m *KernelManager) isAlive(u *User) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/api/system/bootProgress", u.KernelPort))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return http.StatusOK == resp.StatusCode
}

func (m *KernelManager) start(u *User) error {
	workspace := m.workspacePath(u)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return err
	}

	cmd := exec.Command(m.kernelBin, "serve",
		"--workspace", workspace,
		"--port", fmt.Sprintf("%d", u.KernelPort),
		"--lang", m.lang,
	)
	// 网关是唯一入口，内核端口不对外暴露，因此跳过锁屏密码校验
	cmd.Env = append(os.Environ(), "SIYUAN_ACCESS_AUTH_CODE_BYPASS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start kernel for [%s]: %w", u.Name, err)
	}

	m.mu.Lock()
	m.procs[u.Name] = cmd
	m.mu.Unlock()
	go func() {
		_ = cmd.Wait()
		log.Printf("kernel for [%s] exited", u.Name)
	}()

	if err := m.waitBoot(u, 120*time.Second); err != nil {
		return err
	}
	if err := m.ensurePublish(u); err != nil {
		log.Printf("enable publish for [%s] failed: %v", u.Name, err)
	}
	// 首次启动时自动为用户开通内置 S3 同步（若已配置则跳过）
	if err := m.ensureS3Sync(u); err != nil {
		log.Printf("provision s3 sync for [%s] failed: %v", u.Name, err)
	}
	log.Printf("kernel for [%s] is ready on port %d (publish %d)", u.Name, u.KernelPort, u.PublishPort)
	return nil
}

func (m *KernelManager) waitBoot(u *User, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/api/system/bootProgress", u.KernelPort))
		if err == nil {
			var ret struct {
				Code int `json:"code"`
				Data struct {
					Progress int `json:"progress"`
				} `json:"data"`
			}
			err = json.NewDecoder(resp.Body).Decode(&ret)
			resp.Body.Close()
			if err == nil && 0 == ret.Code && 100 <= ret.Data.Progress {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("kernel for [%s] boot timeout", u.Name)
}

// ensurePublish 在用户内核里开启发布服务并固定端口，供 /@user 分享使用
func (m *KernelManager) ensurePublish(u *User) error {
	body, _ := json.Marshal(map[string]any{
		"enable": true,
		"port":   u.PublishPort,
		// 默认拒绝：新用户的文档不会自动公开，须通过「发布访问控制」显式分享
		"defaultDeny": true,
		"auth": map[string]any{
			"enable":   false,
			"accounts": []any{},
		},
	})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("http://127.0.0.1:%d/api/setting/setPublish", u.KernelPort),
		"application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var ret struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return err
	}
	if 0 != ret.Code {
		return fmt.Errorf("setPublish failed: %s", ret.Msg)
	}
	return nil
}

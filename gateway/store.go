// SiYuan Gateway - multi-user front for SiYuan kernels
// Copyright (c) 2026-present
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 网关用户账户，每个用户对应一个独立的工作空间和内核进程
type User struct {
	Name         string    `json:"name"`         // 登录名，同时作为工作空间目录名和分享地址 /@name
	PasswordHash string    `json:"passwordHash"` // bcrypt 哈希
	Created      time.Time `json:"created"`
	KernelPort   int       `json:"kernelPort"`  // 内核 HTTP 端口（仅容器内可达）
	PublishPort  int       `json:"publishPort"` // 发布服务端口（仅容器内可达）
}

// Session 登录会话
type Session struct {
	User    string    `json:"user"`
	Expires time.Time `json:"expires"`
}

// Store 基于 JSON 文件的用户与会话存储，规模 ≤ 数十用户时足够
type Store struct {
	mu       sync.Mutex
	path     string
	Users    map[string]*User    `json:"users"`
	Sessions map[string]*Session `json:"sessions"`
	NextPort int                 `json:"nextPort"` // 下一个可分配的内核端口（每个用户占用两个端口）
}

var usernameRegexp = regexp.MustCompile(`^[a-z][a-z0-9-]{2,31}$`)

var (
	errUserExists   = errors.New("user exists")
	errBadUsername  = errors.New("bad username")
	errBadPassword  = errors.New("bad password")
	errAuthFailed   = errors.New("auth failed")
	errUserNotFound = errors.New("user not found")
)

func loadStore(path string) (*Store, error) {
	s := &Store{
		path:     path,
		Users:    map[string]*User{},
		Sessions: map[string]*Session{},
		NextPort: 30000,
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	if err = json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	if s.NextPort < 30000 {
		s.NextPort = 30000
	}
	return s, nil
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err = os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	if err = os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Register 创建账户并分配端口对
func (s *Store) Register(name, password string) (*User, error) {
	if !usernameRegexp.MatchString(name) {
		return nil, errBadUsername
	}
	if len(password) < 8 {
		return nil, errBadPassword
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Users[name]; ok {
		return nil, errUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &User{
		Name:         name,
		PasswordHash: string(hash),
		Created:      time.Now(),
		KernelPort:   s.NextPort,
		PublishPort:  s.NextPort + 1,
	}
	s.NextPort += 2
	s.Users[name] = u
	return u, s.save()
}

// Authenticate 校验登录名和密码
func (s *Store) Authenticate(name, password string) (*User, error) {
	s.mu.Lock()
	u, ok := s.Users[name]
	s.mu.Unlock()
	if !ok {
		return nil, errAuthFailed
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errAuthFailed
	}
	return u, nil
}

func (s *Store) GetUser(name string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Users[name]
}

// CreateSession 创建会话并持久化
func (s *Store) CreateSession(user string, ttl time.Duration) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gcSessionsLocked()
	s.Sessions[token] = &Session{User: user, Expires: time.Now().Add(ttl)}
	return token, s.save()
}

// SessionUser 根据会话令牌返回用户，过期返回空
func (s *Store) SessionUser(token string) *User {
	if token == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.Sessions[token]
	if !ok || time.Now().After(sess.Expires) {
		return nil
	}
	return s.Users[sess.User]
}

func (s *Store) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Sessions, token)
	_ = s.save()
}

func (s *Store) gcSessionsLocked() {
	now := time.Now()
	for t, sess := range s.Sessions {
		if now.After(sess.Expires) {
			delete(s.Sessions, t)
		}
	}
}

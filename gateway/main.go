// SiYuan Gateway - multi-user front for SiYuan kernels
//
// 网关是多用户 SiYuan 服务的唯一公开入口：
//   - 注册 / 登录（bcrypt + 会话 Cookie）
//   - 每个用户一个独立工作空间和内核进程，按会话反向代理
//   - /@user 通过用户内核的发布服务对外只读分享
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	sessionCookie = "gw_session"
	shareCookie   = "gw_share"

	sessionTTL         = 12 * time.Hour
	sessionTTLRemember = 30 * 24 * time.Hour
)

type Gateway struct {
	store        *Store
	kernels      *KernelManager
	inviteCode   string
	secureCookie bool
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func env(key, def string) string {
	if v := os.Getenv(key); "" != v {
		return v
	}
	return def
}

func main() {
	listenAddr := env("GATEWAY_LISTEN", ":6810")
	dataDir := env("GATEWAY_DATA", "/siyuan/users")
	kernelBin := env("GATEWAY_KERNEL_BIN", "/opt/siyuan/kernel")
	lang := env("GATEWAY_LANG", "ru")
	inviteCode := os.Getenv("GATEWAY_INVITE_CODE")

	store, err := loadStore(filepath.Join(dataDir, "gateway.json"))
	if err != nil {
		log.Fatalf("load store: %v", err)
	}

	gw := &Gateway{
		store:        store,
		kernels:      newKernelManager(kernelBin, dataDir, lang),
		inviteCode:   inviteCode,
		secureCookie: "true" == env("GATEWAY_SECURE_COOKIE", "false"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /gw/login", func(w http.ResponseWriter, r *http.Request) { renderLogin(w, "", http.StatusOK) })
	mux.HandleFunc("POST /gw/login", gw.handleLogin)
	mux.HandleFunc("GET /gw/register", func(w http.ResponseWriter, r *http.Request) {
		renderRegister(w, "", "" != gw.inviteCode, http.StatusOK)
	})
	mux.HandleFunc("POST /gw/register", gw.handleRegister)
	mux.HandleFunc("GET /gw/logout", gw.handleLogout)
	mux.HandleFunc("GET /gw/exit-share", gw.handleExitShare)
	mux.HandleFunc("GET /gw/account", gw.handleAccount)
	mux.HandleFunc("GET /gw/whoami", gw.handleWhoami)
	mux.HandleFunc("/", gw.handleRoot)

	log.Printf("gateway is listening on %s (data dir %s)", listenAddr, dataDir)
	server := &http.Server{Addr: listenAddr, Handler: mux}
	log.Fatal(server.ListenAndServe())
}

// docIDRegexp 文档块 ID 形如 20210808180117-6v0mkxr
var docIDRegexp = regexp.MustCompile(`^\d{14}-[0-9a-z]{7}$`)

// handleRoot 主路由：/@user[/docID] → 固定分享 Cookie；分享 Cookie → 发布服务；会话 → 用户内核
func (gw *Gateway) handleRoot(w http.ResponseWriter, r *http.Request) {
	// 分享入口：/@user 或 /@user/docID，记住目标并跳转，之后所有资源请求都路由到该用户的发布服务
	if strings.HasPrefix(r.URL.Path, "/@") {
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/@"), "/", 2)
		name := parts[0]
		target := gw.store.GetUser(name)
		if nil == target {
			http.NotFound(w, r)
			return
		}
		if err := gw.kernels.Ensure(target); err != nil {
			log.Printf("ensure kernel for shared [%s]: %v", name, err)
			http.Error(w, "Пространство временно недоступно", http.StatusBadGateway)
			return
		}
		gw.setCookie(w, shareCookie, name, int(sessionTTL.Seconds()))
		location := "/"
		if 2 == len(parts) && docIDRegexp.MatchString(parts[1]) {
			// 内核根路由会保留查询参数并按 UA 选择前端包，前端启动时按 ?id= 打开对应文档
			location = "/?id=" + parts[1]
		}
		http.Redirect(w, r, location, http.StatusFound)
		return
	}

	// 查看他人发布内容优先，直到访问 /gw/exit-share 退出
	if c, err := r.Cookie(shareCookie); nil == err && "" != c.Value {
		if target := gw.store.GetUser(c.Value); nil != target {
			gw.proxyTo(w, r, target.PublishPort, target)
			return
		}
		gw.setCookie(w, shareCookie, "", -1)
	}

	user := gw.currentUser(r)
	if nil == user {
		http.Redirect(w, r, "/gw/login", http.StatusFound)
		return
	}
	gw.proxyTo(w, r, user.KernelPort, user)
}

func (gw *Gateway) proxyTo(w http.ResponseWriter, r *http.Request, port int, u *User) {
	if err := gw.kernels.Ensure(u); err != nil {
		log.Printf("ensure kernel for [%s]: %v", u.Name, err)
		http.Error(w, "Не удалось запустить рабочее пространство, попробуйте обновить страницу", http.StatusBadGateway)
		return
	}
	target := &url.URL{Scheme: "http", Host: fmt.Sprintf("127.0.0.1:%d", port)}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy to [%s:%d] failed: %v", u.Name, port, err)
		http.Error(w, "Рабочее пространство не отвечает", http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

func (gw *Gateway) currentUser(r *http.Request) *User {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil
	}
	return gw.store.SessionUser(c.Value)
}

func (gw *Gateway) setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   gw.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

func (gw *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	user, err := gw.store.Authenticate(username, password)
	if err != nil {
		time.Sleep(time.Second) // 简单的暴力破解限速
		renderLogin(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	ttl := sessionTTL
	if "1" == r.FormValue("remember") {
		ttl = sessionTTLRemember
	}
	token, err := gw.store.CreateSession(user.Name, ttl)
	if err != nil {
		renderLogin(w, "Внутренняя ошибка, попробуйте ещё раз", http.StatusInternalServerError)
		return
	}
	gw.setCookie(w, sessionCookie, token, int(ttl.Seconds()))
	gw.setCookie(w, shareCookie, "", -1)

	go func() {
		if err := gw.kernels.Ensure(user); err != nil {
			log.Printf("prestart kernel for [%s]: %v", user.Name, err)
		}
	}()
	http.Redirect(w, r, "/", http.StatusFound)
}

func (gw *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	inviteRequired := "" != gw.inviteCode
	username := strings.TrimSpace(strings.ToLower(r.FormValue("username")))
	password := r.FormValue("password")

	if inviteRequired && gw.inviteCode != strings.TrimSpace(r.FormValue("invite")) {
		renderRegister(w, "Неверный код приглашения", inviteRequired, http.StatusForbidden)
		return
	}
	if password != r.FormValue("password2") {
		renderRegister(w, "Пароли не совпадают", inviteRequired, http.StatusBadRequest)
		return
	}
	if "gw" == username { // 保留给网关自身路由
		renderRegister(w, "Это имя занято", inviteRequired, http.StatusBadRequest)
		return
	}

	user, err := gw.store.Register(username, password)
	switch err {
	case nil:
	case errUserExists:
		renderRegister(w, "Такой пользователь уже существует", inviteRequired, http.StatusConflict)
		return
	case errBadUsername:
		renderRegister(w, "Имя: 3–32 символа, строчные латинские буквы, цифры и дефис, первая — буква", inviteRequired, http.StatusBadRequest)
		return
	case errBadPassword:
		renderRegister(w, "Пароль должен быть не короче 8 символов", inviteRequired, http.StatusBadRequest)
		return
	default:
		log.Printf("register [%s]: %v", username, err)
		renderRegister(w, "Внутренняя ошибка, попробуйте ещё раз", inviteRequired, http.StatusInternalServerError)
		return
	}

	if err = gw.kernels.Ensure(user); err != nil {
		log.Printf("first start kernel for [%s]: %v", user.Name, err)
		renderRegister(w, "Аккаунт создан, но пространство не запустилось — попробуйте войти", inviteRequired, http.StatusBadGateway)
		return
	}

	token, err := gw.store.CreateSession(user.Name, sessionTTL)
	if err != nil {
		http.Redirect(w, r, "/gw/login", http.StatusFound)
		return
	}
	gw.setCookie(w, sessionCookie, token, int(sessionTTL.Seconds()))
	http.Redirect(w, r, "/", http.StatusFound)
}

func (gw *Gateway) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); nil == err {
		gw.store.DeleteSession(c.Value)
	}
	gw.setCookie(w, sessionCookie, "", -1)
	gw.setCookie(w, shareCookie, "", -1)
	http.Redirect(w, r, "/gw/login", http.StatusFound)
}

func (gw *Gateway) handleExitShare(w http.ResponseWriter, r *http.Request) {
	gw.setCookie(w, shareCookie, "", -1)
	http.Redirect(w, r, "/", http.StatusFound)
}

// handleWhoami 返回当前登录名，前端用它拼接分享链接
func (gw *Gateway) handleWhoami(w http.ResponseWriter, r *http.Request) {
	user := gw.currentUser(r)
	if nil == user {
		http.Error(w, `{"user":""}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"user": user.Name})
}

func (gw *Gateway) handleAccount(w http.ResponseWriter, r *http.Request) {
	user := gw.currentUser(r)
	if nil == user {
		http.Redirect(w, r, "/gw/login", http.StatusFound)
		return
	}
	shared := false
	if c, err := r.Cookie(shareCookie); nil == err && "" != c.Value {
		shared = true
	}
	renderAccount(w, user.Name, shared)
}

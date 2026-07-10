// SiYuan Gateway - multi-user front for SiYuan kernels
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package main

import (
	"html/template"
	"net/http"
)

const pageStyle = `
:root { color-scheme: light dark; }
* { box-sizing: border-box; }
body {
  margin: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center;
  font: 16px/1.5 -apple-system, "Segoe UI", Roboto, sans-serif;
  background: light-dark(#f5f6f8, #1e2227); color: light-dark(#1f2329, #d7dde3);
}
.card {
  width: 380px; max-width: calc(100vw - 32px); padding: 32px;
  background: light-dark(#fff, #262b32); border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0,0,0,.12);
}
.card h1 { margin: 0 0 4px; font-size: 22px; }
.card p.sub { margin: 0 0 24px; font-size: 13px; opacity: .65; }
label { display: block; font-size: 13px; margin: 14px 0 4px; opacity: .8; }
input[type=text], input[type=password] {
  width: 100%; padding: 10px 12px; font-size: 15px; border-radius: 8px;
  border: 1px solid light-dark(#d0d5db, #3a414b);
  background: light-dark(#fff, #1e2227); color: inherit;
}
button {
  width: 100%; margin-top: 22px; padding: 11px; font-size: 15px; font-weight: 600;
  border: 0; border-radius: 8px; background: #3575f0; color: #fff; cursor: pointer;
}
button:hover { background: #2b62cc; }
.err { margin: 16px 0 0; padding: 10px 12px; border-radius: 8px; font-size: 14px;
  background: light-dark(#fdecec, #46262a); color: light-dark(#b3261e, #ff8f88); }
.alt { margin-top: 20px; font-size: 14px; text-align: center; }
a { color: #3575f0; text-decoration: none; }
.remember { display: flex; align-items: center; gap: 8px; margin-top: 14px; font-size: 14px; }
.remember input { margin: 0; }
ul.info { padding-left: 18px; font-size: 14px; }
`

var loginTmpl = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="ru"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Вход — SiYuan</title><style>` + pageStyle + `</style></head>
<body><form class="card" method="post" action="/gw/login">
<h1>SiYuan</h1><p class="sub">Вход в личное пространство</p>
<label for="u">Имя пользователя</label>
<input id="u" type="text" name="username" autocomplete="username" autofocus required>
<label for="p">Пароль</label>
<input id="p" type="password" name="password" autocomplete="current-password" required>
<div class="remember"><input id="r" type="checkbox" name="remember" value="1"><label for="r" style="margin:0">Запомнить меня на 30 дней</label></div>
{{if .Error}}<div class="err">{{.Error}}</div>{{end}}
<button type="submit">Войти</button>
<div class="alt">Нет аккаунта? <a href="/gw/register">Зарегистрироваться</a></div>
</form></body></html>`))

var registerTmpl = template.Must(template.New("register").Parse(`<!DOCTYPE html>
<html lang="ru"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Регистрация — SiYuan</title><style>` + pageStyle + `</style></head>
<body><form class="card" method="post" action="/gw/register">
<h1>Регистрация</h1><p class="sub">Будет создано личное рабочее пространство</p>
<label for="u">Имя пользователя</label>
<input id="u" type="text" name="username" autocomplete="username" pattern="[a-z][a-z0-9-]{2,31}" title="3–32 символа: строчные латинские буквы, цифры и дефис; первая — буква" autofocus required>
<label for="p">Пароль (не короче 8 символов)</label>
<input id="p" type="password" name="password" autocomplete="new-password" minlength="8" required>
<label for="p2">Пароль ещё раз</label>
<input id="p2" type="password" name="password2" autocomplete="new-password" minlength="8" required>
{{if .InviteRequired}}
<label for="i">Код приглашения</label>
<input id="i" type="text" name="invite" required>
{{end}}
{{if .Error}}<div class="err">{{.Error}}</div>{{end}}
<button type="submit">Создать аккаунт</button>
<div class="alt">Уже есть аккаунт? <a href="/gw/login">Войти</a></div>
</form></body></html>`))

var accountTmpl = template.Must(template.New("account").Parse(`<!DOCTYPE html>
<html lang="ru"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Аккаунт — SiYuan</title><style>` + pageStyle + `</style></head>
<body><div class="card">
<h1>{{.User}}</h1><p class="sub">Управление аккаунтом</p>
<ul class="info">
<li><a href="/">Открыть моё пространство</a></li>
<li>Публичная ссылка на ваши материалы: <a href="/@{{.User}}">/@{{.User}}</a><br>
<small>Что именно видно другим, настраивается в SiYuan: дерево документов → меню «Ещё» → «Контроль доступа публикации».</small></li>
{{if .Shared}}<li>Сейчас вы просматриваете чужую публикацию — <a href="/gw/exit-share">вернуться в своё пространство</a></li>{{end}}
<li><a href="/gw/logout">Выйти</a></li>
</ul>
</div></body></html>`))

func renderLogin(w http.ResponseWriter, errMsg string, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = loginTmpl.Execute(w, map[string]any{"Error": errMsg})
}

func renderRegister(w http.ResponseWriter, errMsg string, inviteRequired bool, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = registerTmpl.Execute(w, map[string]any{"Error": errMsg, "InviteRequired": inviteRequired})
}

func renderAccount(w http.ResponseWriter, user string, shared bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = accountTmpl.Execute(w, map[string]any{"User": user, "Shared": shared})
}

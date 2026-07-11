#!/usr/bin/env bash
# Мастер установки SiYuan Gateway: спрашивает параметры, генерирует .env и поднимает стек.
# Запуск: bash install.sh   (из папки gateway/ или из корня — скрипт сам перейдёт куда нужно)
set -euo pipefail

cd "$(dirname "$0")"

BOLD=$'\033[1m'; DIM=$'\033[2m'; GRN=$'\033[32m'; YLW=$'\033[33m'; RED=$'\033[31m'; RST=$'\033[0m'
say()  { printf '%s\n' "$*"; }
head() { printf '\n%s%s%s\n' "$BOLD" "$*" "$RST"; }
ok()   { printf '%s✓%s %s\n' "$GRN" "$RST" "$*"; }
warn() { printf '%s!%s %s\n' "$YLW" "$RST" "$*"; }
die()  { printf '%s✗ %s%s\n' "$RED" "$*" "$RST" >&2; exit 1; }

# Спросить значение с дефолтом: ask VAR "Вопрос" "дефолт"
ask() {
  local __var="$1" __q="$2" __def="${3:-}" __ans
  if [ -n "$__def" ]; then
    read -r -p "$(printf '%s %s[%s]%s: ' "$__q" "$DIM" "$__def" "$RST")" __ans || true
  else
    read -r -p "$(printf '%s: ' "$__q")" __ans || true
  fi
  printf -v "$__var" '%s' "${__ans:-$__def}"
}

# Да/нет: ask_yn VAR "Вопрос" "y|n"
ask_yn() {
  local __var="$1" __q="$2" __def="${3:-y}" __ans
  read -r -p "$(printf '%s %s(%s)%s: ' "$__q" "$DIM" "$([ "$__def" = y ] && echo 'Y/n' || echo 'y/N')" "$RST")" __ans || true
  __ans="${__ans:-$__def}"
  case "$__ans" in [Yy]*) printf -v "$__var" 'true' ;; *) printf -v "$__var" 'false' ;; esac
}

gen_secret() {
  if command -v openssl >/dev/null 2>&1; then openssl rand -hex 24
  else head -c 24 /dev/urandom | od -An -tx1 | tr -d ' \n'; fi
}

# ── Проверки окружения ───────────────────────────────────────
head "SiYuan Gateway — установка"
command -v docker >/dev/null 2>&1 || die "Docker не найден. Установите Docker и повторите."
if docker compose version >/dev/null 2>&1; then DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then DC="docker-compose"
else die "Docker Compose не найден (ни 'docker compose', ни 'docker-compose')."; fi
ok "Docker и Compose найдены ($DC)"

# ── Существующий .env ────────────────────────────────────────
REUSE_ENV=false
if [ -f .env ]; then
  ask_yn KEEP "Найден .env. Использовать его без изменений?" y
  [ "$KEEP" = true ] && REUSE_ENV=true
fi

# ── Режим развёртывания ──────────────────────────────────────
head "Режим развёртывания"
say "  1) prod  — тянуть готовый образ из GHCR (рекомендуется для сервера)"
say "  2) local — собрать образ локально из исходников"
ask MODE "Выберите режим (1/2)" "1"
if [ "$MODE" = "2" ]; then
  COMPOSE_ARGS="-f docker-compose.yml -f docker-compose.local.yml"
  say "${DIM}Локальная сборка: docker-compose.yml + docker-compose.local.yml${RST}"
else
  COMPOSE_ARGS="-f docker-compose.prod.yml"
  say "${DIM}Прод: docker-compose.prod.yml (образ ghcr.io/mihazzzold/siyuan-gateway:latest)${RST}"
fi

# ── Сбор параметров ──────────────────────────────────────────
if [ "$REUSE_ENV" = false ]; then
  head "Параметры"
  ask GATEWAY_PORT "Порт шлюза" "6810"
  ask TZ "Часовой пояс" "Europe/Moscow"
  ask GATEWAY_INVITE_CODE "Код приглашения (пусто = открытая регистрация)" ""
  ask_yn BEHIND_PROXY "Шлюз будет за HTTPS-прокси (nginx/caddy)?" y
  GATEWAY_SECURE_COOKIE="$BEHIND_PROXY"
  ask MINIO_ROOT_USER "Пользователь MinIO (S3)" "siyuan"
  ask MINIO_ROOT_PASSWORD "Пароль MinIO (пусто = сгенерировать)" ""
  [ -z "$MINIO_ROOT_PASSWORD" ] && { MINIO_ROOT_PASSWORD="$(gen_secret)"; ok "Пароль MinIO сгенерирован"; }
  ask MINIO_CONSOLE_PORT "Порт веб-консоли MinIO" "9001"
  ask GATEWAY_S3_KEY_SECRET "Секрет ключа шифрования (пусто = сгенерировать)" ""
  [ -z "$GATEWAY_S3_KEY_SECRET" ] && { GATEWAY_S3_KEY_SECRET="$(gen_secret)"; ok "Секрет ключа сгенерирован"; }

  cat > .env <<EOF
GATEWAY_PORT=$GATEWAY_PORT
TZ=$TZ
GATEWAY_INVITE_CODE=$GATEWAY_INVITE_CODE
GATEWAY_SECURE_COOKIE=$GATEWAY_SECURE_COOKIE
MINIO_ROOT_USER=$MINIO_ROOT_USER
MINIO_ROOT_PASSWORD=$MINIO_ROOT_PASSWORD
MINIO_CONSOLE_PORT=$MINIO_CONSOLE_PORT
GATEWAY_S3_KEY_SECRET=$GATEWAY_S3_KEY_SECRET
EOF
  chmod 600 .env
  ok ".env создан (доступ 600)"
  warn "СОХРАНИТЕ GATEWAY_S3_KEY_SECRET и пароль MinIO — при потере зашифрованные S3-данные не восстановить."
fi

# ── Запуск ───────────────────────────────────────────────────
head "Запуск"
ask_yn GO "Поднять стек сейчас?" y
if [ "$GO" = true ]; then
  if [ "$MODE" = "2" ]; then
    $DC $COMPOSE_ARGS --env-file .env up -d --build
  else
    $DC $COMPOSE_ARGS --env-file .env pull
    $DC $COMPOSE_ARGS --env-file .env up -d
  fi
  PORT="$(grep -E '^GATEWAY_PORT=' .env | cut -d= -f2)"
  head "Готово ✓"
  ok  "Шлюз: http://localhost:${PORT:-6810}  (страница входа /gw/login)"
  say "${DIM}Для внешнего доступа поставьте HTTPS-прокси на этот порт (и проксируйте WebSocket /ws).${RST}"
  say "${DIM}Автодеплой: добавьте секрет PORTAINER_WEBHOOK_URL в GitHub — CI будет дёргать его после сборки.${RST}"
else
  say "Стек не запущен. Позже: $DC $COMPOSE_ARGS --env-file .env up -d"
fi

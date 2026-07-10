# SiYuan multi-user gateway image
# 构建上下文必须是仓库根目录：docker build -f gateway/Dockerfile .
FROM --platform=$BUILDPLATFORM node:22 AS node-build

ARG NPM_REGISTRY=

WORKDIR /app
ADD app/package.json app/pnpm* app/.npmrc .

RUN <<EORUN
#!/bin/bash -e
corepack enable
corepack install --global $(node -e 'console.log(require("./package.json").packageManager)')
npm config set registry ${NPM_REGISTRY}
pnpm install --silent
EORUN

ADD app/ .
RUN <<EORUN
#!/bin/bash -e
pnpm run build
node scripts/trimChangelogs.js
mkdir /artifacts
mv appearance stage guide /artifacts/
if [ -d changelogs ]; then mv changelogs /artifacts/; fi
EORUN

FROM golang:1.25-alpine AS go-build

RUN <<EORUN
#!/bin/sh -e
apk add --no-cache gcc musl-dev
go env -w GO111MODULE=on
go env -w CGO_ENABLED=1
EORUN

WORKDIR /kernel
ADD kernel/go.* .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg \
    go mod download

ADD kernel/ .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg \
    go build -tags fts5 -ldflags "-s -w"

WORKDIR /gateway
ADD gateway/ .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg \
    go mod tidy && CGO_ENABLED=0 go build -ldflags "-s -w" -o gateway .

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Europe/Moscow
ENV HOME=/home/siyuan
ENV RUN_IN_CONTAINER=true
ENV GATEWAY_DATA=/siyuan/users
EXPOSE 6810

WORKDIR /opt/siyuan/
COPY --from=go-build --chmod=755 /kernel/kernel .
COPY --from=go-build --chmod=755 /gateway/gateway .
COPY --from=node-build /artifacts .

ENTRYPOINT ["/opt/siyuan/gateway"]

# --- base: 共通の Go ツールチェーン ----------------------------------------
FROM golang:1.25-alpine AS base
RUN apk add --no-cache git ca-certificates
WORKDIR /app

# --- local: ホットリロードの開発サーバー。`make ci` が必要とするツールも同梱 ---
# Laravel比較: Laravel版の docker/app イメージが `composer install --dev` で
# Pint/Larastan を使えるようにしているのと同じ役割。ここでは開発ツール
# (air, migrate, golangci-lint, swag)を、パッケージマネージャー経由ではなく
# Go のバイナリとしてイメージに直接インストールしている。
FROM base AS local
# アプリのソースはホストから実行時にバインドマウントされる(compose.yaml 参照)。
# .git も一緒にマウントされるが、`go build` は通常これを読んで VCS 情報を
# 埋め込もうとする。コンテナ内の root ユーザーはホストがマウントしたその
# .git の所有者ではないため、git がそれを拒否し("detected dubious ownership")、
# ビルド自体が丸ごと失敗する。git の所有者チェックと戦うのではなく、
# ローカル開発では VCS スタンプ自体をスキップする。
ENV GOFLAGS=-buildvcs=false
RUN go install github.com/air-verse/air@v1.61.1 && \
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install golang.org/x/tools/cmd/goimports@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# --- builder: 本番用バイナリのビルド -----------------------------------------
FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

# --- runtime: 最小構成の本番イメージ -----------------------------------------
# Laravel比較: Laravel版が本番相当のローカルモードで動かす app+nginx の組み合わせ
# を置き換える。Go のバイナリ自身が HTTP を処理するので、PHP-FPM/nginx のような
# 構成分割が要らない — distroless イメージに静的バイナリが1つあるだけ。
FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /out/api /api
EXPOSE 8000
ENTRYPOINT ["/api"]

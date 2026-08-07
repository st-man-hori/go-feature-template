# --- base: shared Go toolchain -------------------------------------------
FROM golang:1.25-alpine AS base
RUN apk add --no-cache git ca-certificates
WORKDIR /app

# --- local: hot-reload dev server, with the same tools `make ci` needs ---
# Laravel comparison: this is the equivalent of the Laravel template's
# docker/app image having Pint/Larastan available via `composer install
# --dev`; here the dev tools (air, migrate, golangci-lint, swag) are Go
# binaries installed straight into the image instead of vendored via a
# package manager.
FROM base AS local
RUN go install github.com/air-verse/air@v1.61.1 && \
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install golang.org/x/tools/cmd/goimports@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# --- builder: compiles the production binary ------------------------------
FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

# --- runtime: minimal production image -------------------------------------
# Laravel comparison: replaces the app+nginx pair the Laravel template runs
# in prod-like local mode. A Go binary serves HTTP itself, so there's no
# PHP-FPM/nginx split to configure — one static binary in a distroless image.
FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /out/api /api
EXPOSE 8000
ENTRYPOINT ["/api"]

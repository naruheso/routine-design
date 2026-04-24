# 開発用 Dockerfile — Go + Air（ホットリロード）
FROM golang:1.23-alpine

WORKDIR /app

# ビルドに必要なツール
RUN apk add --no-cache git gcc musl-dev

# Air（ホットリロード）
RUN go install github.com/air-verse/air@v1.61.0

# sqlc（SQLからGoコード自動生成）
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

EXPOSE 8080

# 起動時に依存解決 → Air でホットリロード開始
CMD ["sh", "-c", "go mod tidy && air -c .air.toml"]

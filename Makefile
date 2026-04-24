.PHONY: up down build logs ps clean db-reset

# === 開発環境 ===

## 起動（全サービス）
up:
	docker compose up -d

## 起動（ログ表示付き）
up-log:
	docker compose up

## 停止
down:
	docker compose down

## リビルド＆起動
build:
	docker compose up -d --build

## ログ表示
logs:
	docker compose logs -f

## 各サービスのログ
logs-api:
	docker compose logs -f backend

logs-web:
	docker compose logs -f frontend

logs-db:
	docker compose logs -f postgres

## コンテナ状態確認
ps:
	docker compose ps

# === データベース ===

## DB初期化（データ全削除＆再作成）
db-reset:
	docker compose down -v
	docker compose up -d postgres
	@echo "DB初期化完了"

## psql接続
db-connect:
	docker compose exec postgres psql -U postgres -d routine_design

# === バックエンド ===

## sqlc コード生成
sqlc:
	docker compose exec backend sqlc generate

## Go テスト実行
test-backend:
	docker compose exec backend go test ./... -v

# === フロントエンド ===

## フロントエンド テスト実行
test-frontend:
	docker compose exec frontend npm run test

# === クリーンアップ ===

## 全削除（ボリューム含む）
clean:
	docker compose down -v --rmi local
	@echo "クリーンアップ完了"

# 開発用 Dockerfile — Vue 3 + Vite（ホットリロード）
FROM node:20-alpine

WORKDIR /app

EXPOSE 5173

# 起動時に依存インストール → Vite dev server 起動
CMD ["sh", "-c", "npm install && npm run dev -- --host 0.0.0.0"]

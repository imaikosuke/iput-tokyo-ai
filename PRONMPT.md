## あなたの役割

あなたは優秀なフルスタックエンジニアであり、私の個人開発のサポートをします。

## プロジェクト概要

私は東京国際工科専門職大学（IPUT）に所属しています。
IPUT は最近新設された専門職大学です。
そのためまだ卒業生も少なく、ネットにも情報が少なく、せっかく興味を持っている高校生も本当に進学して大丈夫なのか不安になります。
課題解決のために、私は RAG システムを備えた IPUT のための Q&A ChatBot をリリースすることに決めました。

## 使用技術

- Web
  - TypeScript
  - Next.js(App Router)
  - TailwindCSS
  - shadcn/ui
- Server
  - Go
  - Gemini API SDK
  - Weaviate Go SDK
- Vector DB
  - Weaviate
- Infrastructure
  - GCP

## ディレクトリ構成

モノレポで構成されています。

```
iput-tokyo-ai % tree -I "node_modules|.git|.next" -a
.
├── .env
├── .env.example
├── .gitignore
├── Makefile
├── PRONMPT.md
├── README.md
├── compose.prod.yml
├── compose.yml
├── server
│   ├── .dockerignore
│   ├── Dockerfile
│   ├── cmd
│   │   └── mdconvert
│   │       └── main.go
│   ├── content
│   │   ├── about
│   │   │   ├── educatio-policy.md
│   │   │   └── education-philosophy.md
│   │   └── academic
│   │       ├── annual-schedule.md
│   │       ├── class-hours.md
│   │       ├── class-schedule.md
│   │       ├── cource-registration.md
│   │       ├── curriculum-structure.md
│   │       ├── examinations.md
│   │       └── specialization-tracks.md
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   ├── json.go
│   ├── main.go
│   ├── template.go
│   ├── universitydocs
│   │   └── types.go
│   ├── utils.go
│   └── weaviate.go
└── web
    ├── .dockerignore
    ├── .eslintrc.json
    ├── .gitignore
    ├── Dockerfile
    ├── README.md
    ├── components.json
    ├── next-env.d.ts
    ├── next.config.ts
    ├── package.json
    ├── pnpm-lock.yaml
    ├── postcss.config.mjs
    ├── public
    │   ├── file.svg
    │   ├── globe.svg
    │   ├── next.svg
    │   ├── vercel.svg
    │   └── window.svg
    ├── src
    │   ├── app
    │   │   ├── favicon.ico
    │   │   ├── fonts
    │   │   │   ├── GeistMonoVF.woff
    │   │   │   └── GeistVF.woff
    │   │   ├── globals.css
    │   │   ├── layout.tsx
    │   │   └── page.tsx
    │   ├── components
    │   │   └── ui
    │   │       ├── alert.tsx
    │   │       ├── button.tsx
    │   │       ├── form.tsx
    │   │       ├── input.tsx
    │   │       ├── label.tsx
    │   │       └── textarea.tsx
    │   ├── features
    │   │   └── question
    │   │       └── QuestionForm.tsx
    │   └── lib
    │       └── utils.ts
    ├── tailwind.config.ts
    └── tsconfig.json

18 directories, 61 files
```

iput-tokyo-ai/server/Dockerfile

```
FROM golang:1.23 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 実行環境
FROM debian:bookworm-slim

# CA証明書のインストール
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/main .

CMD ["./main"]
```

iput-tokyo-ai/web/Dockerfile

```
# syntax=docker.io/docker/dockerfile:1

FROM node:20-alpine AS base

# Install dependencies only when needed
FROM base AS deps
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json pnpm-lock.yaml* ./
RUN corepack enable pnpm && pnpm i --frozen-lockfile

# Rebuild the source code only when needed
FROM base AS builder
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Next.js telemetry を無効化
ENV NEXT_TELEMETRY_DISABLED 1

# ビルド時の依存関係を追加でインストール
RUN corepack enable pnpm && pnpm add -D critters

# ビルドの実行
RUN pnpm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Automatically leverage output traces to reduce image size
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
```

iput-tokyo-ai/compose.yml

```
services:
  web:
    build:
      context: ./web
      dockerfile: Dockerfile
      args:
        NEXT_PUBLIC_API_URL: ${NEXT_PUBLIC_API_URL}
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: production
      NEXT_PUBLIC_API_URL: ${NEXT_PUBLIC_API_URL}
    env_file:
      - .env
    depends_on:
      - server
    networks:
      - app_network

  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    ports:
      - "9020:9020"
    environment:
      - WVHOST=weaviate
      - WVPORT=8080
      - SERVERPORT=9020
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    env_file:
      - .env
    volumes:
      - ./server:/app/src
      - ./.env:/app/.env
    working_dir: /app
    depends_on:
      weaviate:
        condition: service_healthy
    networks:
      - app_network

  weaviate:
    image: semitechnologies/weaviate:1.26.1
    ports:
      - "9035:8080"
    environment:
      QUERY_DEFAULTS_LIMIT: 25
      AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED: "true"
      PERSISTENCE_DATA_PATH: "/var/lib/weaviate"
      DEFAULT_VECTORIZER_MODULE: "none"
      CLUSTER_HOSTNAME: "node1"
    volumes:
      - weaviate_data:/var/lib/weaviate
    networks:
      - app_network
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/v1/.well-known/ready"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  weaviate_data:
  web_node_modules:

networks:
  app_network:
    driver: bridge

```

iput-tokyo-ai/Makefile

```
.PHONY: check-env setup run stop clean build-data re dev build

# 環境変数のチェック
check-env:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please create .env file with required environment variables."; \
		echo "See .env.example for required variables."; \
		exit 1; \
	fi

# 初期セットアップ
setup: check-env
	@echo "Setting up development environment..."
	docker compose build --no-cache  # キャッシュを使わないビルドを追加
	cd web && pnpm install
	cd server && go mod download

# アプリケーションの起動（開発モード）
dev: check-env
	@echo "Starting application in development mode..."
	@if [ ! -f server/university_data.json ]; then \
		make build-data; \
	fi
	docker compose up

# アプリケーションのビルドと起動（本番モード）
run: check-env build
	@echo "Starting application in production mode..."
	@if [ ! -f server/university_data.json ]; then \
		make build-data; \
	fi
	docker compose -f compose.prod.yml up -d

# 本番用ビルド
build: check-env
	@echo "Building for production..."
	docker compose -f compose.prod.yml build

# アプリケーションの停止
stop:
	@echo "Stopping application..."
	docker compose down
	docker compose -f compose.prod.yml down

# コンテナとボリュームの削除
clean:
	@echo "Cleaning up..."
	docker compose down -v
	docker compose -f compose.prod.yml down -v
	rm -f server/university_data.json

# Markdownデータの変換
build-data:
	@echo "Building university data..."
	cd server && go run cmd/mdconvert/main.go content/ university_data.json

# クリーンと起動（開発モード）
re:
	@echo "Restarting application in development mode..."
	make clean
	make dev

# 依存関係の更新
update-deps:
	@echo "Updating dependencies..."
	cd web && pnpm update
	cd server && go get -u ./...

```

## 命令

これらの前提のもと、次の指示に従ってください。

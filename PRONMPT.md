## あなたの役割
あなたは優秀なフルスタックエンジニアであり、私の個人開発のサポートをします。

## プロジェクト概要
私は東京国際工科専門職大学（IPUT）に所属しています。
IPUTは最近新設された専門職大学です。
そのためまだ卒業生も少なく、ネットにも情報が少なく、せっかく興味を持っている高校生も本当に進学して大丈夫なのか不安になります。
課題解決のために、私はRAGシステムを備えたIPUTのためのQ&A ChatBotをリリースすることに決めました。

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
モノリポで構成されています。
```
iput-tokyo-ai % tree -I "node_modules|.git" -a
.
├── .env
├── .gitignore
├── Makefile
├── PRONMPT.md
├── README.md
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
│   ├── university_data.json
│   ├── universitydocs
│   │   └── types.go
│   ├── utils.go
│   └── weaviate.go
└── web
    ├── .dockerignore
    ├── .eslintrc.json
    ├── .gitignore
    ├── .next
    │   ├── app-build-manifest.json
    │   ├── build-manifest.json
    │   ├── cache
    │   │   ├── config.json
    │   │   └── webpack
    │   │       └── client-development
    │   │           ├── 0.pack.gz
    │   │           ├── index.pack.gz
    │   │           └── index.pack.gz.old
    │   ├── package.json
    │   ├── react-loadable-manifest.json
    │   ├── server
    │   │   ├── app-paths-manifest.json
    │   │   ├── interception-route-rewrite-manifest.js
    │   │   ├── middleware-build-manifest.js
    │   │   ├── middleware-manifest.json
    │   │   ├── middleware-react-loadable-manifest.js
    │   │   ├── next-font-manifest.js
    │   │   ├── next-font-manifest.json
    │   │   ├── pages-manifest.json
    │   │   ├── server-reference-manifest.js
    │   │   └── server-reference-manifest.json
    │   ├── static
    │   │   ├── chunks
    │   │   │   └── polyfills.js
    │   │   └── development
    │   │       ├── _buildManifest.js
    │   │       └── _ssgManifest.js
    │   ├── trace
    │   └── types
    │       ├── cache-life.d.ts
    │       └── package.json
    ├── Dockerfile
    ├── README.md
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
    │   └── app
    │       ├── favicon.ico
    │       ├── fonts
    │       │   ├── GeistMonoVF.woff
    │       │   └── GeistVF.woff
    │       ├── globals.css
    │       ├── layout.tsx
    │       └── page.tsx
    ├── tailwind.config.ts
    └── tsconfig.json

22 directories, 75 files
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

WORKDIR /app
COPY --from=builder /build/main .

CMD ["./main"]
```

iput-tokyo-ai/web/Dockerfile
```
FROM node:20-slim

WORKDIR /app

COPY package*.json ./
COPY pnpm-lock.yaml ./

RUN npm install -g pnpm
RUN pnpm install

COPY . .

CMD ["pnpm", "dev"]
```

iput-tokyo-ai/compose.yml
```
services:
  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    volumes:
      - ./web:/app
      - /app/node_modules
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:9020
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

networks:
  app_network:
    driver: bridge

```

iput-tokyo-ai/Makefile
```
.PHONY: check-env setup run stop clean build-data re

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
	docker compose build
	cd web && pnpm install
	cd server && go mod download

# アプリケーションの起動
run: check-env
	@echo "Starting application..."
	@if [ ! -f server/university_data.json ]; then \
		make build-data; \
	fi
	docker compose up

# アプリケーションの停止
stop:
	@echo "Stopping application..."
	docker compose down

# コンテナとボリュームの削除
clean:
	@echo "Cleaning up..."
	docker compose down -v
	rm -f server/university_data.json

# Markdownデータの変換
build-data:
	@echo "Building university data..."
	cd server && go run cmd/mdconvert/main.go content/ university_data.json

# クリーンと起動
re:
	@echo "Restarting application..."
	make clean
	make run


```

## 命令
これらの前提のもと、次の指示に従ってください。

"use client";

import { useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function Home() {
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!question.trim()) return;

    setIsLoading(true);
    setError(null);
    setAnswer("");

    try {
      const response = await fetch("http://localhost:9020/query/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ content: question }),
      });

      if (!response.ok) {
        throw new Error(`APIエラー: ${response.status}`);
      }

      const data = await response.text();
      setAnswer(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "予期せぬエラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="container mx-auto p-4 max-w-3xl">
      <h1 className="text-2xl font-bold mb-6">東京国際工科専門職大学 Q&A</h1>

      <form onSubmit={handleSubmit} className="space-y-4 mb-6">
        <div>
          <label htmlFor="question" className="block text-sm font-medium mb-2">
            質問を入力してください
          </label>
          <textarea
            id="question"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            className="w-full p-2 border rounded-md min-h-[100px]"
            placeholder="例: 情報工学科について教えてください"
          />
        </div>
        <button
          type="submit"
          disabled={isLoading || !question.trim()}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-blue-300 disabled:cursor-not-allowed"
        >
          {isLoading ? "送信中..." : "送信"}
        </button>
      </form>

      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {answer && (
        <div className="border rounded-md p-4 bg-gray-50">
          <h2 className="font-medium mb-2">回答:</h2>
          <div className="whitespace-pre-wrap">{answer}</div>
        </div>
      )}

      {isLoading && (
        <div className="flex justify-center items-center space-x-2">
          <div className="animate-spin h-5 w-5 border-2 border-blue-500 rounded-full border-t-transparent"></div>
          <span>回答を生成中...</span>
        </div>
      )}
    </main>
  );
}

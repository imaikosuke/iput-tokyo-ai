.PHONY: check-env setup run stop clean build-data re dev build rebuild copy-files

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

# サーバーの再ビルドと起動
rebuild:
	@echo "Rebuilding server and starting application..."
	docker compose down
	docker compose build --no-cache server
	make dev

# 依存関係の更新
update-deps:
	@echo "Updating dependencies..."
	cd web && pnpm update
	cd server && go get -u ./...

copy-files:
	@if [ -z "$(dir)" ]; then \
		echo "Error: ディレクトリを指定してください (例: make copy-files dir=server/cmd)"; \
		exit 1; \
	fi
	@./copyfiles $(dir)

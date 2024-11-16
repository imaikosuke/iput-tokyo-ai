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

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
# # サーバー起動前にWeaviateが完全に起動するのを待つスクリプトを追加
# COPY <<'EOF' /wait-for-weaviate.sh
# #!/bin/sh
# set -e

# until wget --spider -q http://weaviate:8080/v1/.well-known/ready; do
#   echo 'Waiting for Weaviate...'
#   sleep 1
# done

# echo 'Weaviate is up - executing command'
# EOF

# RUN chmod +x /wait-for-weaviate.sh

# # Weaviateが起動するのを待ってから、Markdownの変換とサーバー起動を行う
# CMD ["/bin/sh", "-c", "/wait-for-weaviate.sh && go run cmd/mdconvert/main.go content/ university_data.json && go run main.go json.go weaviate.go"]

# ホットリロード用のair をインストール
# RUN go install github.com/air-verse/air@latest

# airの設定ファイルをコピー
# COPY .air.toml ./

# airを使用して実行
# ENTRYPOINT ["air", "-c", ".air.toml"]


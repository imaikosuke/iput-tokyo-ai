package main

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"google.golang.org/api/option"
)

const GENERATIVE_MODEL = "gemini-1.5-flash"
const EMBEDDING_MODEL = "text-embedding-004"

type ragServer struct {
	ctx      context.Context        // コンテキスト
	wvClient *weaviate.Client       // Weaviateクライアント
	genModel *genai.GenerativeModel // GenerativeAIモデル
	embModel *genai.EmbeddingModel  // EmbeddingAIモデル
}

// CORSミドルウェアの設定
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Printf("Warning: .env file not found")
	}
	log.Print("env: ", os.Getenv("GEMINI_API_KEY"))

	// Weaviateクライアントの初期化
	ctx := context.Background()
	wvClient, err := initWeaviate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// GenerativeAIクライアントの初期化
	apiKey := os.Getenv("GEMINI_API_KEY")
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer genaiClient.Close()

	// サーバーの初期化
	server := &ragServer{
		ctx:      ctx,
		wvClient: wvClient,
		genModel: genaiClient.GenerativeModel(GENERATIVE_MODEL),
		embModel: genaiClient.EmbeddingModel(EMBEDDING_MODEL),
	}

	// APIエンドポイントの設定
	mux := http.NewServeMux()
	mux.HandleFunc("POST /add/", server.addDocumentsHandler)
	mux.HandleFunc("POST /query/", server.queryHandler)

	// CORSミドルウェアの適用
	handler := corsMiddleware(mux)

	// サーバーの起動
	port := cmp.Or(os.Getenv("SERVERPORT"), "9020")
	address := ":" + port
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, handler))
}

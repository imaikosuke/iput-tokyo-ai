package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"

	"github.com/google/generative-ai-go/genai"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func (rs *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Starting addDocumentsHandler")
	type document struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Category   string   `json:"category"`
		Tags       []string `json:"tags"`
		Department string   `json:"department"`
		UpdatedAt  string   `json:"updated_at"`
	}
	type addRequest struct {
		Documents []document `json:"documents"`
	}
	addRequestDocuments := &addRequest{}

	err := readRequestJSON(req, addRequestDocuments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// チャンカーの設定を構築
	cfg, err := config.NewConfigBuilder().
		WithMaxTokens(512).
		WithMinTokens(100).
		WithOverlapTokens(50).
		WithJapaneseConfig(config.NewDefaultJapaneseConfig()).
		Build()

	if err != nil {
		http.Error(w, fmt.Sprintf("configuring chunker: %v", err), http.StatusInternalServerError)
		return
	}

	// チャンカーの初期化
	chunker, err := chunking.NewChunker(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("initializing chunker: %v", err), http.StatusInternalServerError)
		return
	}

	var allObjects []*models.Object

	// ドキュメントごとの処理
	for i, doc := range addRequestDocuments.Documents {
		log.Printf("Processing document %d: %s", i, doc.Title)
		log.Printf("Document content length: %d", len(doc.Content))

		// コンテンツをチャンクに分割
		chunks, err := chunker.ChunkDocument(doc.Content)
		if err != nil {
			log.Printf("Error chunking document: %v", err)
			http.Error(w, fmt.Sprintf("chunking document: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Document '%s' was split into %d chunks", doc.Title, len(chunks))
		for i, chunk := range chunks {
			log.Printf("Chunk %d: %d tokens, %d-%d chars",
				i, chunk.TokenCount, chunk.StartChar, chunk.EndChar)
		}

		// チャンクごとのembedding用バッチを作成
		batch := rs.embModel.NewBatch()
		for _, chunk := range chunks {
			fullText := fmt.Sprintf(
				"Title: %s\nCategory: %s\nDepartment: %s\nContent: %s",
				doc.Title, doc.Category, doc.Department, chunk.Content,
			)
			batch.AddContent(genai.Text(fullText))
		}

		// バッチembedding処理
		rsp, err := rs.embModel.BatchEmbedContents(rs.ctx, batch)
		if err != nil {
			http.Error(w, fmt.Sprintf("batch embedding: %v", err), http.StatusInternalServerError)
			return
		}

		// チャンクごとにWeaviateオブジェクトを作成
		for i, chunk := range chunks {
			obj := &models.Object{
				Class: "Document",
				Properties: map[string]any{
					"title":       doc.Title,
					"content":     chunk.Content,
					"category":    doc.Category,
					"tags":        doc.Tags,
					"department":  doc.Department,
					"updatedAt":   doc.UpdatedAt,
					"chunkIndex":  chunk.Index,
					"totalChunks": len(chunks),
					"startChar":   chunk.StartChar,
					"endChar":     chunk.EndChar,
					"tokenCount":  chunk.TokenCount,
					"precedence":  chunk.Precedence,
				},
				Vector: rsp.Embeddings[i].Values,
			}
			allObjects = append(allObjects, obj)
		}
	}

	// Weaviateへの保存
	log.Printf("storing %v objects in weaviate", len(allObjects))
	_, err = rs.wvClient.Batch().ObjectsBatcher().WithObjects(allObjects...).Do(rs.ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("storing in weaviate: %v", err), http.StatusInternalServerError)
		return
	}

	renderJSON(w, map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %d document chunks", len(allObjects)),
	})
}

type Response struct {
	Answer string `json:"answer"`
}

func (rs *ragServer) queryHandler(w http.ResponseWriter, req *http.Request) {
	type queryRequest struct {
		Content string
	}
	qr := &queryRequest{}
	err := readRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// クエリの埋め込み処理
	rsp, err := rs.embModel.EmbedContent(rs.ctx, genai.Text(qr.Content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Weaviateでの類似検索（上位3チャンクを取得）
	gql := rs.wvClient.GraphQL()
	result, err := gql.Get().
		WithClassName("Document").
		WithFields(
			graphql.Field{Name: "title"},
			graphql.Field{Name: "content"},
			graphql.Field{Name: "category"},
			graphql.Field{Name: "department"},
			graphql.Field{Name: "chunkIndex"},
			graphql.Field{Name: "totalChunks"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{
				{Name: "certainty"},
			}},
		).
		WithNearVector(
			gql.NearVectorArgBuilder().
				WithVector(rsp.Embedding.Values).
				WithCertainty(0.7)).
		WithLimit(5).
		Do(rs.ctx)

	log.Printf("Query response: %+v", result.Data)
	if werr := combinedWeaviateError(result, err); werr != nil {
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	contents, err := decodeGetResults(result)
	if err != nil {
		http.Error(w, fmt.Errorf("reading weaviate response: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved %d relevant chunks from Weaviate", len(contents))

	// RAGクエリの生成と実行
	ragQuery := fmt.Sprintf(GetRAGTemplate(), qr.Content, strings.Join(contents, "\n\n---\n\n"))
	log.Printf("RAG query:\n%s", ragQuery)
	resp, err := rs.genModel.GenerateContent(rs.ctx, genai.Text(ragQuery))
	if err != nil {
		log.Printf("calling generative model: %v", err.Error())
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	if len(resp.Candidates) != 1 {
		log.Printf("got %v candidates, expected 1", len(resp.Candidates))
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	var respTexts []string
	for _, part := range resp.Candidates[0].Content.Parts {
		if pt, ok := part.(genai.Text); ok {
			respTexts = append(respTexts, string(pt))
		} else {
			log.Printf("bad type of part: %v", pt)
			http.Error(w, "generative model error", http.StatusInternalServerError)
			return
		}
	}

	renderJSON(w, Response{Answer: strings.Join(respTexts, "\n")})
}

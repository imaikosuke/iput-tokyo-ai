package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func (rs *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
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

	// バッチembedding処理
	batch := rs.embModel.NewBatch()
	for _, doc := range addRequestDocuments.Documents {
		fullText := fmt.Sprintf("Title: %s\nCategory: %s\nDepartment: %s\nContent: %s",
			doc.Title, doc.Category, doc.Department, doc.Content)
		batch.AddContent(genai.Text(fullText))
	}

	log.Printf("invoking embedding model with %v documents", len(addRequestDocuments.Documents))
	rsp, err := rs.embModel.BatchEmbedContents(rs.ctx, batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rsp.Embeddings) != len(addRequestDocuments.Documents) {
		http.Error(w, "embedded batch size mismatch", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated embeddings for %d documents, vector length: %d",
		len(rsp.Embeddings),
		len(rsp.Embeddings[0].Values))

	// Weaviateオブジェクトの作成
	objects := make([]*models.Object, len(addRequestDocuments.Documents))
	for i, doc := range addRequestDocuments.Documents {
		objects[i] = &models.Object{
			Class: "Document",
			Properties: map[string]any{
				"title":      doc.Title,
				"content":    doc.Content,
				"category":   doc.Category,
				"tags":       doc.Tags,
				"department": doc.Department,
				"updatedAt":  doc.UpdatedAt,
			},
			Vector: rsp.Embeddings[i].Values,
		}
	}

	// Weaviateへの保存
	log.Printf("storing %v objects in weaviate", len(objects))
	_, err = rs.wvClient.Batch().ObjectsBatcher().WithObjects(objects...).Do(rs.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %d documents", len(addRequestDocuments.Documents)),
	})
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

	// Weaviateでの類似検索
	gql := rs.wvClient.GraphQL()
	result, err := gql.Get().
		WithClassName("Document").
		WithFields(
			graphql.Field{Name: "title"},
			graphql.Field{Name: "content"},
			graphql.Field{Name: "category"},
			graphql.Field{Name: "department"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{
				{Name: "certainty"},
			}},
		).
		WithNearVector(
			gql.NearVectorArgBuilder().
				WithVector(rsp.Embedding.Values).
				WithCertainty(0.5)).
		WithLimit(1).
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
	log.Printf("Retrieved %d relevant documents from Weaviate", len(contents))

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

	renderJSON(w, strings.Join(respTexts, "\n"))
}

package universitydocs

import "github.com/weaviate/weaviate/entities/models"

type Document struct {
	// ID         string   `json:"id" yaml:"id"`
	Title      string   `json:"title" yaml:"title"`
	Content    string   `json:"content" yaml:"content"`
	Category   string   `json:"category" yaml:"category"`
	Tags       []string `json:"tags" yaml:"tags"`
	Department string   `json:"department" yaml:"department"`
	UpdatedAt  string   `json:"updated_at" yaml:"updated_at"`
}

type AddDocumentsRequest struct {
	Documents []Document `json:"documents"`
}

// Weaviateへの保存用にドキュメントを変換する関数
func ConvertToWeaviateObject(doc Document) *models.Object {
	return &models.Object{
		Class: "Document",
		Properties: map[string]any{
			// "id":         doc.ID,
			"title":      doc.Title,
			"content":    doc.Content,
			"category":   doc.Category,
			"tags":       doc.Tags,
			"department": doc.Department,
			"updatedAt":  doc.UpdatedAt,
		},
	}
}

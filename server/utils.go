package main

import (
	"fmt"
	"log"

	"github.com/weaviate/weaviate/entities/models"
)

// Weaviate GraphQLのレスポンスをデコードし、ドキュメントの内容を含む文字列のリストを返す
func decodeGetResults(result *models.GraphQLResponse) ([]string, error) {
	data, ok := result.Data["Get"]
	if !ok {
		return nil, fmt.Errorf("don't have get key in response")
	}
	document, ok := data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid get key in response")
	}
	slices, ok := document["Document"].([]any)
	if !ok {
		return nil, fmt.Errorf("document is not a list of results")
	}

	var out []string
	for index, slice := range slices {
		slicedData, ok := slice.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid element in list of documents")
		}

		title, _ := slicedData["title"].(string)
		content, _ := slicedData["content"].(string)
		category, _ := slicedData["category"].(string)
		department, _ := slicedData["department"].(string)

		additional, _ := slicedData["_additional"].(map[string]any)
		certainty := 0.0
		if additional != nil {
			if cert, ok := additional["certainty"].(float64); ok {
				certainty = cert
			}
		}

		log.Printf("Document %d: %s (certainty: %.3f)", index+1, title, certainty)

		docContent := fmt.Sprintf("タイトル: %s\nカテゴリ: %s\n所属: %s\n\n%s",
			title, category, department, content)

		out = append(out, docContent)
	}
	return out, nil
}

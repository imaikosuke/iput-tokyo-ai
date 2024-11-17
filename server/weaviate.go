package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func initWeaviate(ctx context.Context) (*weaviate.Client, error) {
	host := cmp.Or(os.Getenv("WVHOST"), "weaviate")
	port := cmp.Or(os.Getenv("WVPORT"), "8080")

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   fmt.Sprintf("%s:%s", host, port),
		Scheme: "http",
	})
	if err != nil {
		return nil, fmt.Errorf("initializing weaviate: %w", err)
	}

	// Weaviateのスキーマを初期化
	cls := &models.Class{
		Class:      "Document",
		Vectorizer: "none",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"string"},
			},
			{
				Name:     "content",
				DataType: []string{"text"},
			},
			{
				Name:     "category",
				DataType: []string{"string"},
			},
			{
				Name:     "tags",
				DataType: []string{"string[]"},
			},
			{
				Name:     "department",
				DataType: []string{"string"},
			},
			{
				Name:     "updatedAt",
				DataType: []string{"string"},
			},
			// チャンク関連の新しいプロパティ
			{
				Name:     "chunkIndex",
				DataType: []string{"int"},
			},
			{
				Name:     "totalChunks",
				DataType: []string{"int"},
			},
			{
				Name:     "startChar",
				DataType: []string{"int"},
			},
			{
				Name:     "endChar",
				DataType: []string{"int"},
			},
			{
				Name:     "tokenCount",
				DataType: []string{"int"},
			},
		},
	}

	// Weaviateのクラスが存在しない場合は作成
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		exists, err := client.Schema().ClassExistenceChecker().WithClassName(cls.Class).Do(ctx)
		if err == nil {
			if !exists {
				err = client.Schema().ClassCreator().WithClass(cls).Do(ctx)
				if err != nil {
					return nil, fmt.Errorf("creating weaviate class: %w", err)
				}
			}
			return client, nil
		}

		log.Printf("Failed to connect to Weaviate (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("failed to initialize Weaviate after %d attempts", maxRetries)
}

func combinedWeaviateError(result *models.GraphQLResponse, err error) error {
	if err != nil {
		return err
	}
	if len(result.Errors) != 0 {
		var ss []string
		for _, e := range result.Errors {
			ss = append(ss, e.Message)
		}
		return fmt.Errorf("weaviate error: %v", ss)
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Document は university-data-structure で定義した構造体と同じ
type Document struct {
	// ID         string   `json:"id" yaml:"id"`
	Title      string   `json:"title" yaml:"title"`
	Content    string   `json:"content" yaml:"content"`
	Category   string   `json:"category" yaml:"category"`
	Tags       []string `json:"tags" yaml:"tags"`
	Department string   `json:"department" yaml:"department"`
	UpdatedAt  string   `json:"updated_at" yaml:"updated_at"`
}

// フロントマターを含むMarkdownファイルの構造
type MarkdownFile struct {
	Frontmatter string
	Content     string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: mdconvert <input-directory> <output-file>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputFile := os.Args[2]

	docs, err := processDirectory(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}

	output := struct {
		Documents []Document `json:"documents"`
	}{
		Documents: docs,
	}

	jsonData, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}
}

// 引数の`dir`で指定されたディレクトリ内のMarkdownファイルを処理して`Document`のスライスを返す
func processDirectory(dir string) ([]Document, error) {
	var docs []Document

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			doc, err := processFile(path)
			if err != nil {
				return fmt.Errorf("processing %s: %w", path, err)
			}
			docs = append(docs, doc)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return docs, nil
}

// 引数の`filename`で指定されたMarkdownファイルを処理して`Document`を返す
func processFile(filename string) (Document, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Document{}, err
	}

	parts := strings.Split(string(content), "---")
	if len(parts) < 3 {
		return Document{}, fmt.Errorf("invalid markdown format in %s", filename)
	}

	var doc Document
	err = yaml.Unmarshal([]byte(parts[1]), &doc)
	if err != nil {
		return Document{}, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// IDが指定されていない場合はファイル名をIDとして使用
	// if doc.ID == "" {
	// 	doc.ID = strings.TrimSuffix(filepath.Base(filename), ".md")
	// }

	// 更新日が指定されていない場合は現在時刻を使用
	if doc.UpdatedAt == "" {
		doc.UpdatedAt = time.Now().Format("2006-01-02")
	}

	// Markdownコンテンツの取得（3番目のパート以降を結合）
	doc.Content = strings.TrimSpace(strings.Join(parts[2:], "---"))

	return doc, nil
}

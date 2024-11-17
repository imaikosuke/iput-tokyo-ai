// server/cmd/chunking/main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
)

// ChunkResult は1つのチャンクの結果を表示するための構造体
type ChunkResult struct {
	Index      int               `json:"index"`
	Content    string            `json:"content"`
	StartChar  int               `json:"startChar"`
	EndChar    int               `json:"endChar"`
	TokenCount int               `json:"tokenCount"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Precedence int               `json:"precedence"`
	References []string          `json:"references,omitempty"`
}

func main() {
	// 基本設定のフラグ
	var (
		maxTokens          = flag.Int("max", 512, "Maximum tokens per chunk")
		minTokens          = flag.Int("min", 100, "Minimum tokens per chunk")
		overlapTokens      = flag.Int("overlap", 50, "Number of overlapping tokens")
		paragraphSeparator = flag.String("sep", "\n\n", "Paragraph separator")
		inputFile          = flag.String("in", "", "Input file path (optional)")
		outputFile         = flag.String("out", "", "Output JSON file path (optional)")
	)

	// 日本語処理の設定フラグ
	var (
		keyParticleWeight = flag.Float64("particle-weight", 1.2, "Weight for key particles in Japanese text")
		topicMarkerWeight = flag.Float64("topic-weight", 1.5, "Weight for topic markers in Japanese text")
	)

	// 文書構造の設定フラグ
	var (
		listItemWeight  = flag.Float64("list-weight", 0.8, "Weight for list items")
		codeBlockWeight = flag.Float64("code-weight", 1.2, "Weight for code blocks")
		tableWeight     = flag.Float64("table-weight", 1.5, "Weight for tables")
	)

	// 出力形式の設定
	var (
		verbose        = flag.Bool("verbose", false, "Show detailed information")
		showMetadata   = flag.Bool("meta", true, "Include metadata in output")
		showReferences = flag.Bool("refs", true, "Include references in output")
	)

	flag.Parse()

	// 設定の構築
	japaneseConfig := &config.JapaneseConfig{
		SentenceEndings:   []string{"。", "！", "？", "…"},
		Brackets:          []string{"（）", "「」", "『』", "［］"},
		KeyParticleWeight: *keyParticleWeight,
		TopicMarkerWeight: *topicMarkerWeight,
	}

	// 設定のビルドとエラーハンドリング
	cfg, err := config.NewConfigBuilder().
		WithMaxTokens(*maxTokens).
		WithMinTokens(*minTokens).
		WithOverlapTokens(*overlapTokens).
		WithParagraphSeparator(*paragraphSeparator).
		WithJapaneseConfig(japaneseConfig).
		WithListItemWeight(*listItemWeight).
		WithCodeBlockWeight(*codeBlockWeight).
		WithTableWeight(*tableWeight).
		Build()

	if err != nil {
		log.Fatalf("Failed to build configuration: %v", err)
	}

	// チャンカーの作成
	chunker, err := chunking.NewChunker(cfg)
	if err != nil {
		log.Fatalf("Failed to create chunker: %v", err)
	}
	// 入力テキストの取得
	inputText := getSampleOrFileContent(*inputFile)

	// チャンク分割の実行
	chunks, err := chunker.ChunkDocument(inputText)
	if err != nil {
		log.Fatalf("Failed to chunk document: %v", err)
	}

	// 結果の表示用に構造体に変換
	results := convertChunksToResults(chunks, *showMetadata, *showReferences)

	// 結果の出力
	if *outputFile != "" {
		outputResults(results, *outputFile)
	} else {
		displayResults(results, chunker.GetConfig(), *verbose)
	}
}

// convertChunksToResults はChunkをChunkResultに変換する
func convertChunksToResults(chunks []models.Chunk, showMetadata, showReferences bool) []ChunkResult {
	results := make([]ChunkResult, len(chunks))
	for i, chunk := range chunks {
		results[i] = ChunkResult{
			Index:      chunk.Index,
			Content:    chunk.Content,
			StartChar:  chunk.StartChar,
			EndChar:    chunk.EndChar,
			TokenCount: chunk.TokenCount,
			Precedence: chunk.Precedence,
		}
		if showMetadata {
			results[i].Metadata = chunk.Metadata
		}
		if showReferences {
			results[i].References = chunk.References
		}
	}
	return results
}

// displayResults は結果をコンソールに表示する
func displayResults(results []ChunkResult, cfg *config.ChunkConfig, verbose bool) {
	fmt.Println("Configuration:")
	displayConfig(cfg)

	fmt.Printf("\nFound %d chunks:\n\n", len(results))
	for _, result := range results {
		displayChunk(result, verbose)
	}
}

// displayConfig は設定情報を表示する
func displayConfig(cfg *config.ChunkConfig) {
	fmt.Printf("Basic Settings:\n")
	fmt.Printf("- Max Tokens: %d\n", cfg.MaxTokens)
	fmt.Printf("- Min Tokens: %d\n", cfg.MinTokens)
	fmt.Printf("- Overlap Tokens: %d\n", cfg.OverlapTokens)

	if cfg.JapaneseConfig != nil {
		fmt.Printf("\nJapanese Processing:\n")
		fmt.Printf("- Key Particle Weight: %.2f\n", cfg.JapaneseConfig.KeyParticleWeight)
		fmt.Printf("- Topic Marker Weight: %.2f\n", cfg.JapaneseConfig.TopicMarkerWeight)
	}

	fmt.Printf("\nStructure Weights:\n")
	fmt.Printf("- List Items: %.2f\n", cfg.ListItemWeight)
	fmt.Printf("- Code Blocks: %.2f\n", cfg.CodeBlockWeight)
	fmt.Printf("- Tables: %.2f\n", cfg.TableWeight)
}

// 他のヘルパー関数は変更なし
func getSampleOrFileContent(inputFile string) string {
	if inputFile != "" {
		content, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Failed to read input file: %v", err)
		}
		return string(content)
	}
	return getSampleText()
}

func outputResults(results []ChunkResult, outputFile string) {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal results to JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}
	fmt.Printf("Results written to %s\n", outputFile)
}

// displayChunk は個々のチャンク情報を表示する
func displayChunk(result ChunkResult, verbose bool) {
	fmt.Printf("Chunk %d:\n", result.Index+1)
	fmt.Printf("- Tokens: %d\n", result.TokenCount)
	fmt.Printf("- Position: %d-%d\n", result.StartChar, result.EndChar)
	if verbose {
		fmt.Printf("- Precedence: %d\n", result.Precedence)
		if len(result.References) > 0 {
			fmt.Printf("- References: %s\n", strings.Join(result.References, ", "))
		}
		if len(result.Metadata) > 0 {
			fmt.Printf("- Metadata:\n")
			for k, v := range result.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
	}
	fmt.Printf("Content:\n%s\n", strings.TrimSpace(result.Content))
	fmt.Println(strings.Repeat("-", 80))
}

// getSampleText はサンプルテキストを返す
func getSampleText() string {
	return `---
title: "東京国際工科専門職大学について"
category: "大学案内"
tags: ["概要", "教育理念", "カリキュラム"]
department: "全学部共通"
updated_at: "2024-03-01"
---

# 東京国際工科専門職大学について

東京国際工科専門職大学（IPUT）は、最先端のテクノロジーと実践的な専門教育を提供する専門職大学です。

## 教育理念

私たちは、技術革新と創造性を重視し、グローバルな視点を持つIT人材の育成に力を入れています。
産業界との密接な連携により、実践的なスキルと理論的知識の両方を身につけることができます。

## カリキュラム特徴

- プロジェクトベースの学習
- 第一線で活躍する実務家教員による指導
- 充実した英語教育プログラム

### 実践的な学び

1年次から実践的なプロジェクトに参加し、実際の課題解決に取り組みます。
企業との共同プロジェクトも多数実施しています。`
}

// pkg/chunking/utils/text.go
package utils

import (
	"strings"
)

// TextProcessor はテキスト処理のためのユーティリティ構造体
type TextProcessor struct{}

// NewTextProcessor は新しいTextProcessorインスタンスを作成
func NewTextProcessor() *TextProcessor {
	return &TextProcessor{}
}

// NormalizeWhitespace は空白文字を正規化
func (tp *TextProcessor) NormalizeWhitespace(text string) string {
	// 連続する空白を単一の空白に置換
	normalized := strings.Join(strings.Fields(text), " ")
	return normalized
}

// CountWords はテキスト内の単語数を計算
func (tp *TextProcessor) CountWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

// TruncateText は指定された長さでテキストを切り詰める
func (tp *TextProcessor) TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// 文の途中で切らないように調整
	lastSpace := strings.LastIndex(text[:maxLength], " ")
	if lastSpace == -1 {
		return text[:maxLength]
	}
	return text[:lastSpace] + "..."
}

// ExtractHeading は見出しレベルとテキストを抽出
func (tp *TextProcessor) ExtractHeading(line string) (int, string) {
	if !strings.HasPrefix(line, "#") {
		return 0, line
	}

	level := 0
	for i, char := range line {
		if char != '#' {
			return level, strings.TrimSpace(line[i:])
		}
		level++
	}
	return level, ""
}

// SplitParagraphs は段落単位でテキストを分割
func (tp *TextProcessor) SplitParagraphs(text string) []string {
	// 改行を正規化
	normalized := strings.ReplaceAll(text, "\r\n", "\n")

	// 連続する改行で分割
	paragraphs := strings.Split(normalized, "\n\n")

	// 空の段落を除去
	var result []string
	for _, p := range paragraphs {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

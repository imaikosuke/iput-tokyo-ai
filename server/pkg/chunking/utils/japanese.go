// pkg/chunking/utils/japanese.go
package utils

import (
	"strings"
	"unicode"
)

// JapaneseProcessor は日本語テキスト処理のためのユーティリティ構造体
type JapaneseProcessor struct {
	// 句読点などの区切り文字
	Separators []string
	// 助詞リスト
	Particles []string
	// トピックマーカー
	TopicMarkers []string
}

// NewJapaneseProcessor は新しいJapaneseProcessorインスタンスを作成
func NewJapaneseProcessor() *JapaneseProcessor {
	return &JapaneseProcessor{
		Separators:   []string{"。", "、", "！", "？", "…"},
		Particles:    []string{"は", "が", "を", "に", "へ", "と", "で", "から", "まで", "より"},
		TopicMarkers: []string{"は", "が", "について", "に関して"},
	}
}

// NormalizeJapaneseText は日本語テキストを正規化
func (jp *JapaneseProcessor) NormalizeJapaneseText(text string) string {
	// 全角スペースを半角スペースに変換
	text = strings.ReplaceAll(text, "　", " ")

	// 改行を正規化
	text = strings.ReplaceAll(text, "\r\n", "\n")

	return text
}

// SplitJapaneseSentences は日本語の文を分割
func (jp *JapaneseProcessor) SplitJapaneseSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		current.WriteRune(runes[i])

		// 文末かどうかをチェック
		if i+1 < len(runes) && jp.isSentenceEnd(string(runes[i])) {
			sentences = append(sentences, strings.TrimSpace(current.String()))
			current.Reset()
		}
	}

	// 残りの文字列を追加
	if current.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(current.String()))
	}

	return sentences
}

// CountJapaneseTokens は日本語テキストのトークン数を計算
func (jp *JapaneseProcessor) CountJapaneseTokens(text string) int {
	text = jp.NormalizeJapaneseText(text)
	var tokenCount int

	// 文字ごとに処理
	for _, r := range text {
		if jp.isJapaneseCharacter(r) {
			tokenCount++
		} else if !unicode.IsSpace(r) {
			// 英数字などは1文字を1トークンとしてカウント
			tokenCount++
		}
	}

	return tokenCount
}

// IsParticle は文字列が助詞かどうかを判定
func (jp *JapaneseProcessor) IsParticle(word string) bool {
	for _, particle := range jp.Particles {
		if word == particle {
			return true
		}
	}
	return false
}

// IsTopicMarker は文字列がトピックマーカーかどうかを判定
func (jp *JapaneseProcessor) IsTopicMarker(word string) bool {
	for _, marker := range jp.TopicMarkers {
		if word == marker {
			return true
		}
	}
	return false
}

// isSentenceEnd は文字が文末表現かどうかを判定
func (jp *JapaneseProcessor) isSentenceEnd(char string) bool {
	for _, sep := range jp.Separators {
		if char == sep {
			return true
		}
	}
	return false
}

// isJapaneseCharacter は文字が日本語文字かどうかを判定
func (jp *JapaneseProcessor) isJapaneseCharacter(r rune) bool {
	return unicode.In(r, unicode.Hiragana) ||
		unicode.In(r, unicode.Katakana) ||
		unicode.In(r, unicode.Han)
}

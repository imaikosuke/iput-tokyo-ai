// pkg/chunking/config/japanese.go
package config

// JapaneseConfig は日本語テキスト処理のための設定
type JapaneseConfig struct {
	// 文の区切り文字
	SentenceEndings []string `json:"sentenceEndings"` // 文末表現のリスト

	// 特殊な区切り文字
	Brackets []string `json:"brackets"` // 括弧のペアリスト

	// 重要度の重み付け
	KeyParticleWeight float64 `json:"keyParticleWeight"` // 重要な助詞の重み
	TopicMarkerWeight float64 `json:"topicMarkerWeight"` // トピックマーカーの重み
}

// NewDefaultJapaneseConfig は日本語設定のデフォルト値を返す
func NewDefaultJapaneseConfig() *JapaneseConfig {
	return &JapaneseConfig{
		SentenceEndings:   []string{"。", "！", "？", "…"},
		Brackets:          []string{"（）", "「」", "『』", "［］"},
		KeyParticleWeight: 1.2,
		TopicMarkerWeight: 1.5,
	}
}

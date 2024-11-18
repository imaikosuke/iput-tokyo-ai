// pkg/chunking/config/config.go
package config

// ChunkConfig はチャンキングの設定を保持する構造体
type ChunkConfig struct {
	// 基本設定
	MaxTokens          int    `json:"maxTokens"`          // チャンクの最大トークン数
	MinTokens          int    `json:"minTokens"`          // チャンクの最小トークン数
	OverlapTokens      int    `json:"overlapTokens"`      // チャンク間のオーバーラップトークン数
	ParagraphSeparator string `json:"paragraphSeparator"` // 段落の区切り文字

	// セクション階層の設定
	MaxSectionDepth int       `json:"maxSectionDepth"` // 処理する最大の見出しレベル（0は無制限）
	SectionWeights  []float64 `json:"sectionWeights"`  // 各見出しレベルの重み付け

	// チャンク結合の設定
	MergeThreshold   float64 `json:"mergeThreshold"`   // チャンク結合を行う類似度の閾値
	MaxMergedTokens  int     `json:"maxMergedTokens"`  // 結合後の最大トークン数
	PreserveSections bool    `json:"preserveSections"` // セクション境界を維持するか

	// メタデータ処理の設定
	MetadataFields []string `json:"metadataFields"` // 抽出するメタデータフィールド
	TagSeparator   string   `json:"tagSeparator"`   // タグの区切り文字

	// 文書構造の解析設定
	ListItemWeight  float64 `json:"listItemWeight"`  // リスト項目の重み付け
	CodeBlockWeight float64 `json:"codeBlockWeight"` // コードブロックの重み付け
	TableWeight     float64 `json:"tableWeight"`     // テーブルの重み付け

	// 優先度計算の設定
	BasePrecedence int     `json:"basePrecedence"` // 基本優先度スコア
	ContentFactor  float64 `json:"contentFactor"`  // 内容長による優先度の重み
	DepthFactor    float64 `json:"depthFactor"`    // 階層深さによる優先度の重み

	// 言語固有の設定
	JapaneseConfig *JapaneseConfig `json:"japaneseConfig"` // 日本語固有の設定
}

// NewDefaultConfig は推奨されるデフォルトのチャンク設定を返す
func NewDefaultConfig() *ChunkConfig {
	return &ChunkConfig{
		// 基本設定
		MaxTokens:          512,
		MinTokens:          100,
		OverlapTokens:      50,
		ParagraphSeparator: "\n\n",

		// セクション階層の設定
		MaxSectionDepth: 6,
		SectionWeights:  []float64{1.0, 0.8, 0.6, 0.4, 0.3, 0.2},

		// チャンク結合の設定
		MergeThreshold:   0.7,
		MaxMergedTokens:  768,
		PreserveSections: true,

		// メタデータ処理の設定
		MetadataFields: []string{"title", "category", "tags", "department", "updated_at"},
		TagSeparator:   ",",

		// 文書構造の解析設定
		ListItemWeight:  0.8,
		CodeBlockWeight: 1.2,
		TableWeight:     1.5,

		// 優先度計算の設定
		BasePrecedence: 100,
		ContentFactor:  0.01,
		DepthFactor:    -0.2,

		// 日本語固有の設定を追加
		JapaneseConfig: NewDefaultJapaneseConfig(),
	}
}

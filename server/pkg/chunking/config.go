package chunking

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

// JapaneseConfig は日本語テキスト処理のための設定
type JapaneseConfig struct {
	// 文の区切り文字
	SentenceEndings []string `json:"sentenceEndings"` // 文末表現のリスト（「。」「！」「？」など）

	// 特殊な区切り文字
	Brackets []string `json:"brackets"` // 括弧のペアリスト

	// 重要度の重み付け
	KeyParticleWeight float64 `json:"keyParticleWeight"` // 重要な助詞の重み
	TopicMarkerWeight float64 `json:"topicMarkerWeight"` // トピックマーカー（「は」「が」など）の重み
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

		// 日本語固有の設定
		JapaneseConfig: &JapaneseConfig{
			SentenceEndings:   []string{"。", "！", "？", "…"},
			Brackets:          []string{"（）", "「」", "『』", "［］"},
			KeyParticleWeight: 1.2,
			TopicMarkerWeight: 1.5,
		},
	}
}

// WithMaxTokens は最大トークン数を設定する
func (c *ChunkConfig) WithMaxTokens(maxTokens int) *ChunkConfig {
	c.MaxTokens = maxTokens
	return c
}

// WithMinTokens は最小トークン数を設定する
func (c *ChunkConfig) WithMinTokens(minTokens int) *ChunkConfig {
	c.MinTokens = minTokens
	return c
}

// WithOverlapTokens はオーバーラップトークン数を設定する
func (c *ChunkConfig) WithOverlapTokens(overlapTokens int) *ChunkConfig {
	c.OverlapTokens = overlapTokens
	return c
}

// WithParagraphSeparator は段落区切り文字を設定する
func (c *ChunkConfig) WithParagraphSeparator(separator string) *ChunkConfig {
	c.ParagraphSeparator = separator
	return c
}

// WithJapaneseConfig は日本語設定を更新する
func (c *ChunkConfig) WithJapaneseConfig(config *JapaneseConfig) *ChunkConfig {
	c.JapaneseConfig = config
	return c
}

// WithMetadataFields はメタデータフィールドを設定する
func (c *ChunkConfig) WithMetadataFields(fields []string) *ChunkConfig {
	c.MetadataFields = fields
	return c
}

// WithSectionWeights はセクションの重み付けを設定する
func (c *ChunkConfig) WithSectionWeights(weights []float64) *ChunkConfig {
	c.SectionWeights = weights
	return c
}

// WithListItemWeight はリスト項目の重み付けを設定する
func (c *ChunkConfig) WithListItemWeight(weight float64) *ChunkConfig {
	c.ListItemWeight = weight
	return c
}

// WithCodeBlockWeight はコードブロックの重み付けを設定する
func (c *ChunkConfig) WithCodeBlockWeight(weight float64) *ChunkConfig {
	c.CodeBlockWeight = weight
	return c
}

// WithTableWeight はテーブルの重み付けを設定する
func (c *ChunkConfig) WithTableWeight(weight float64) *ChunkConfig {
	c.TableWeight = weight
	return c
}

// WithMergeThreshold はチャンク結合の閾値を設定する
func (c *ChunkConfig) WithMergeThreshold(threshold float64) *ChunkConfig {
	c.MergeThreshold = threshold
	return c
}

// WithMaxMergedTokens は結合後の最大トークン数を設定する
func (c *ChunkConfig) WithMaxMergedTokens(maxTokens int) *ChunkConfig {
	c.MaxMergedTokens = maxTokens
	return c
}

// WithPreserveSections はセクション境界の維持設定を更新する
func (c *ChunkConfig) WithPreserveSections(preserve bool) *ChunkConfig {
	c.PreserveSections = preserve
	return c
}

// Validate は設定の妥当性を検証する
func (c *ChunkConfig) Validate() error {
	// バリデーションロジックの実装
	return nil
}

// Clone は設定のディープコピーを作成する
func (c *ChunkConfig) Clone() *ChunkConfig {
	newConfig := &ChunkConfig{
		MaxTokens:          c.MaxTokens,
		MinTokens:          c.MinTokens,
		OverlapTokens:      c.OverlapTokens,
		ParagraphSeparator: c.ParagraphSeparator,
		MaxSectionDepth:    c.MaxSectionDepth,
		SectionWeights:     make([]float64, len(c.SectionWeights)),
		MergeThreshold:     c.MergeThreshold,
		MaxMergedTokens:    c.MaxMergedTokens,
		PreserveSections:   c.PreserveSections,
		MetadataFields:     make([]string, len(c.MetadataFields)),
		TagSeparator:       c.TagSeparator,
		ListItemWeight:     c.ListItemWeight,
		CodeBlockWeight:    c.CodeBlockWeight,
		TableWeight:        c.TableWeight,
		BasePrecedence:     c.BasePrecedence,
		ContentFactor:      c.ContentFactor,
		DepthFactor:        c.DepthFactor,
	}

	copy(newConfig.SectionWeights, c.SectionWeights)
	copy(newConfig.MetadataFields, c.MetadataFields)

	if c.JapaneseConfig != nil {
		newConfig.JapaneseConfig = &JapaneseConfig{
			SentenceEndings:   make([]string, len(c.JapaneseConfig.SentenceEndings)),
			Brackets:          make([]string, len(c.JapaneseConfig.Brackets)),
			KeyParticleWeight: c.JapaneseConfig.KeyParticleWeight,
			TopicMarkerWeight: c.JapaneseConfig.TopicMarkerWeight,
		}
		copy(newConfig.JapaneseConfig.SentenceEndings, c.JapaneseConfig.SentenceEndings)
		copy(newConfig.JapaneseConfig.Brackets, c.JapaneseConfig.Brackets)
	}

	return newConfig
}

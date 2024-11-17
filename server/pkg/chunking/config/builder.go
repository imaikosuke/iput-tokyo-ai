// pkg/chunking/config/builder.go
package config

// ConfigBuilder は設定を構築するためのビルダーパターンを実装する
type ConfigBuilder struct {
	config *ChunkConfig
}

// NewConfigBuilder は新しいConfigBuilderインスタンスを作成する
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: NewDefaultConfig(),
	}
}

// WithMaxTokens は最大トークン数を設定する
func (b *ConfigBuilder) WithMaxTokens(maxTokens int) *ConfigBuilder {
	b.config.MaxTokens = maxTokens
	return b
}

// WithMinTokens は最小トークン数を設定する
func (b *ConfigBuilder) WithMinTokens(minTokens int) *ConfigBuilder {
	b.config.MinTokens = minTokens
	return b
}

// WithOverlapTokens はオーバーラップトークン数を設定する
func (b *ConfigBuilder) WithOverlapTokens(overlapTokens int) *ConfigBuilder {
	b.config.OverlapTokens = overlapTokens
	return b
}

// WithParagraphSeparator は段落区切り文字を設定する
func (b *ConfigBuilder) WithParagraphSeparator(separator string) *ConfigBuilder {
	b.config.ParagraphSeparator = separator
	return b
}

// WithJapaneseConfig は日本語設定を設定する
func (b *ConfigBuilder) WithJapaneseConfig(japaneseConfig *JapaneseConfig) *ConfigBuilder {
	b.config.JapaneseConfig = japaneseConfig
	return b
}

// WithListItemWeight はリスト項目の重み付けを設定する
func (b *ConfigBuilder) WithListItemWeight(weight float64) *ConfigBuilder {
	b.config.ListItemWeight = weight
	return b
}

// WithCodeBlockWeight はコードブロックの重み付けを設定する
func (b *ConfigBuilder) WithCodeBlockWeight(weight float64) *ConfigBuilder {
	b.config.CodeBlockWeight = weight
	return b
}

// WithTableWeight はテーブルの重み付けを設定する
func (b *ConfigBuilder) WithTableWeight(weight float64) *ConfigBuilder {
	b.config.TableWeight = weight
	return b
}

// WithSectionWeights はセクションの重み付けを設定する
func (b *ConfigBuilder) WithSectionWeights(weights []float64) *ConfigBuilder {
	b.config.SectionWeights = weights
	return b
}

// WithMetadataFields はメタデータフィールドを設定する
func (b *ConfigBuilder) WithMetadataFields(fields []string) *ConfigBuilder {
	b.config.MetadataFields = fields
	return b
}

// WithPreserveSections はセクション境界の維持設定を更新する
func (b *ConfigBuilder) WithPreserveSections(preserve bool) *ConfigBuilder {
	b.config.PreserveSections = preserve
	return b
}

// Build は設定を構築する
func (b *ConfigBuilder) Build() (*ChunkConfig, error) {
	if err := b.config.Validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}

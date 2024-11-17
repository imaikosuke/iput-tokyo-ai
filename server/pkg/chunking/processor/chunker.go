// pkg/chunking/processor/chunker.go
package processor

import (
	"strings"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
)

// ContentChunker はコンテンツのチャンキングを担当
type ContentChunker struct {
	config *config.ChunkConfig
}

// NewContentChunker は新しいContentChunkerを作成
func NewContentChunker(cfg *config.ChunkConfig) *ContentChunker {
	return &ContentChunker{
		config: cfg,
	}
}

// Chunk はコンテンツをチャンクに分割
func (c *ContentChunker) Chunk(content string) ([]models.Chunk, error) {
	sections := c.parseSections(content)
	var chunks []models.Chunk
	var headingStack []string

	currentPosition := 0

	for _, section := range sections {
		// 見出し階層の更新
		headingStack = c.updateHeadingStack(headingStack, section.Level, section.Title)

		// セクションの分割
		sectionChunks := c.chunkSection(section.Content)

		// 各チャンクの位置情報とメタデータを設定
		for i := range sectionChunks {
			chunk := &sectionChunks[i]
			if i == 0 {
				chunk.StartChar = currentPosition
			} else {
				chunk.StartChar = chunks[len(chunks)-1].EndChar + 1
			}
			chunk.EndChar = chunk.StartChar + len(chunk.Content)

			// 参照情報の設定
			chunk.References = make([]string, len(headingStack))
			copy(chunk.References, headingStack)

			// 優先度の計算
			chunk.Precedence = c.calculatePrecedence(section.Level, len(chunk.Content))

			chunks = append(chunks, *chunk)
		}

		currentPosition = chunks[len(chunks)-1].EndChar + 1
	}

	return chunks, nil
}

// parseSections はマークダウンからセクション階層を解析する
func (c *ContentChunker) parseSections(content string) []models.Section {
	var sections []models.Section
	lines := strings.Split(content, "\n")
	var currentSection *models.Section
	var buffer strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			// 前のセクションを保存
			if currentSection != nil && buffer.Len() > 0 {
				currentSection.Content = buffer.String()
				sections = append(sections, *currentSection)
				buffer.Reset()
			}

			// 新しいセクションの開始
			level := strings.Count(line, "#")
			title := strings.TrimSpace(strings.TrimLeft(line, "# "))
			currentSection = &models.Section{
				Title: title,
				Level: level,
			}
		}
		buffer.WriteString(line + "\n")
	}

	// 最後のセクションを保存
	if currentSection != nil && buffer.Len() > 0 {
		currentSection.Content = buffer.String()
		sections = append(sections, *currentSection)
	}

	return sections
}

// updateHeadingStack は見出し階層を更新する
func (c *ContentChunker) updateHeadingStack(stack []string, level int, title string) []string {
	// レベルが1未満の場合は新しいスタックを作成
	if level < 1 {
		return []string{title}
	}

	// 現在のスタックの深さを確認
	currentDepth := len(stack)

	// 新しい見出しが現在の階層より深い場合
	if level > currentDepth {
		return append(stack, title)
	}

	// 新しい見出しが現在の階層と同じか浅い場合
	newStack := make([]string, level-1)
	if level-1 > 0 {
		copy(newStack, stack[:level-1])
	}
	return append(newStack, title)
}

// chunkSection はセクションの内容をチャンクに分割
func (c *ContentChunker) chunkSection(content string) []models.Chunk {
	if content = strings.TrimSpace(content); content == "" {
		return nil
	}

	var chunks []models.Chunk
	paragraphs := strings.Split(content, c.config.ParagraphSeparator)
	var currentChunk strings.Builder
	currentTokens := 0

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		tokens := len(strings.Fields(para)) // 簡易的なトークン計算
		if currentTokens+tokens > c.config.MaxTokens && currentChunk.Len() > 0 {
			// 現在のチャンクを保存
			chunks = append(chunks, models.Chunk{
				Content:    currentChunk.String(),
				TokenCount: currentTokens,
			})
			currentChunk.Reset()
			currentTokens = 0
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentTokens += tokens
	}

	// 残りのコンテンツをチャンクとして追加
	if currentChunk.Len() > 0 {
		chunks = append(chunks, models.Chunk{
			Content:    currentChunk.String(),
			TokenCount: currentTokens,
		})
	}

	return chunks
}

// calculatePrecedence はセクションの重要度を計算する
func (c *ContentChunker) calculatePrecedence(level int, contentLength int) int {
	// レベルが低いほど（より上位の見出しほど）優先度は高く
	// コンテンツが長いほど優先度は高く設定
	basePrecedence := 100 - (level * 20)
	contentFactor := contentLength / 100
	return basePrecedence + contentFactor
}

// pkg/chunking/processor/chunker.go
package processor

import (
	"log"
	"strings"
	"unicode"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/utils"
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
			currentSection = models.NewSection(title, level, "")
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
	if level < 1 {
		return []string{title}
	}

	currentDepth := len(stack)

	if level > currentDepth {
		return append(stack, title)
	}

	newStack := make([]string, level-1)
	if level-1 > 0 {
		copy(newStack, stack[:level-1])
	}
	return append(newStack, title)
}

// chunkSection はセクションの内容をチャンクに分割
func (c *ContentChunker) chunkSection(content string) []models.Chunk {
	// デバッグログの追加
	log.Printf("Starting chunkSection with content length: %d", len(content))

	content = strings.TrimSpace(content)
	if content == "" {
		log.Printf("Empty content after trimming")
		return []models.Chunk{}
	}

	// 入力の最大長をチェック
	if len(content) > c.config.MaxTokens*4 {
		log.Printf("Content exceeds maximum length: %d", len(content))
		runes := []rune(content)
		if len(runes) > c.config.MaxTokens {
			content = string(runes[:c.config.MaxTokens])
			log.Printf("Content truncated to length: %d", len(content))
		}
	}

	var chunks []models.Chunk
	paragraphs := strings.Split(content, c.config.ParagraphSeparator)
	log.Printf("Split into %d paragraphs", len(paragraphs))

	var currentChunk strings.Builder
	var currentTokens int
	japaneseProcessor := utils.NewJapaneseProcessor()

	// チャンクの追加用ヘルパー関数
	addChunk := func(content string, tokens int) {
		if content != "" {
			log.Printf("Adding chunk with length: %d, tokens: %d", len(content), tokens)
			chunk := models.NewChunk(strings.TrimSpace(content))
			chunk.TokenCount = tokens
			chunks = append(chunks, *chunk)
		}
	}

	for i, para := range paragraphs {
		log.Printf("Processing paragraph %d of length: %d", i, len(para))
		para = strings.TrimSpace(para)
		if para == "" {
			log.Printf("Empty paragraph, skipping")
			continue
		}

		// パラグラフの長さチェックを追加
		if len(para) > c.config.MaxTokens*4 {
			log.Printf("Paragraph exceeds maximum length: %d", len(para))
			continue
		}

		// トークン数を計算（日本語処理を使用）
		paraTokens := japaneseProcessor.CountJapaneseTokens(para)
		log.Printf("Paragraph tokens: %d", paraTokens)

		// 特殊な文書要素の重み付けを適用
		weight := 1.0
		if strings.HasPrefix(para, "- ") || strings.HasPrefix(para, "* ") {
			weight = c.config.ListItemWeight
		} else if strings.HasPrefix(para, "```") {
			weight = c.config.CodeBlockWeight
		} else if strings.Contains(para, "|") && strings.Contains(para, "-+-") {
			weight = c.config.TableWeight
		}

		adjustedTokens := int(float64(paraTokens) * weight)
		log.Printf("Adjusted tokens: %d", adjustedTokens)

		// 見出し行の処理
		if strings.HasPrefix(para, "#") {
			if currentChunk.Len() > 0 {
				addChunk(currentChunk.String(), currentTokens)
				currentChunk.Reset()
				currentTokens = 0
			}
			currentChunk.WriteString(para)
			currentTokens = adjustedTokens
			continue
		}

		// チャンクサイズの判断
		if currentTokens+adjustedTokens > c.config.MaxTokens && currentChunk.Len() > 0 {
			log.Printf("Chunk size limit reached. Current: %d, Adding: %d, Max: %d",
				currentTokens, adjustedTokens, c.config.MaxTokens)
			addChunk(currentChunk.String(), currentTokens)
			currentChunk.Reset()
			currentTokens = 0
		}

		// 段落の追加
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentTokens += adjustedTokens

		// 日本語の文末表現による分割
		if c.config.JapaneseConfig != nil && len(para) > 0 {
			sentences := japaneseProcessor.SplitJapaneseSentences(para)
			log.Printf("Split into %d sentences", len(sentences))

			if len(sentences) > 1 && currentTokens >= c.config.MinTokens {
				runesPara := []rune(para)
				if len(runesPara) > 0 {
					lastChar := runesPara[len(runesPara)-1]
					if japaneseProcessor.IsSentenceEnd(string(lastChar)) {
						log.Printf("Sentence end detected: %s", string(lastChar))
						addChunk(currentChunk.String(), currentTokens)
						currentChunk.Reset()
						currentTokens = 0
					}
				}
			}
		}
	}

	// 残りのコンテンツを追加
	if currentChunk.Len() > 0 {
		log.Printf("Adding final chunk")
		addChunk(currentChunk.String(), currentTokens)
	}

	log.Printf("Finished chunking. Created %d chunks", len(chunks))
	return chunks
}

// countTokens は日本語と英語のトークン数を適切に計算
func (c *ContentChunker) countTokens(text string, japaneseProcessor *utils.JapaneseProcessor) int {
	// 全角スペースを半角スペースに変換
	text = strings.ReplaceAll(text, "　", " ")

	// 記号による分割
	text = strings.NewReplacer(
		"、", " 、 ",
		"。", " 。 ",
		"（", " （ ",
		"）", " ） ",
		"「", " 「 ",
		"」", " 」 ",
		"『", " 『 ",
		"』", " 』 ",
		"［", " ［ ",
		"］", " ］ ",
		"#", " # ",
		"-", " - ",
	).Replace(text)

	words := strings.Fields(text)
	total := 0

	for _, word := range words {
		if japaneseProcessor.IsJapaneseCharacter([]rune(word)[0]) {
			// 日本語文字列は文字単位でカウント
			charCount := 0
			for _, r := range word {
				if unicode.IsSpace(r) {
					continue
				}
				charCount++
			}
			total += charCount

			// 助詞やトピックマーカーの重み付け
			if c.config.JapaneseConfig != nil {
				if japaneseProcessor.IsParticle(word) {
					total = int(float64(total) * c.config.JapaneseConfig.KeyParticleWeight)
				}
				if japaneseProcessor.IsTopicMarker(word) {
					total = int(float64(total) * c.config.JapaneseConfig.TopicMarkerWeight)
				}
			}
		} else {
			// 英数字などはワード単位でカウント
			total++
		}
	}

	return total
}

// calculatePrecedence はセクションの重要度を計算する
func (c *ContentChunker) calculatePrecedence(level int, contentLength int) int {
	basePrecedence := 100 - (level * 20)
	contentFactor := contentLength / 100
	return basePrecedence + contentFactor
}

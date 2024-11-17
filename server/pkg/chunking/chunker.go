package chunking

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Chunk は単一のチャンクを表す構造体
type Chunk struct {
	Content    string            `json:"content"`    // チャンクの内容
	StartChar  int               `json:"startChar"`  // 元のテキストでの開始位置
	EndChar    int               `json:"endChar"`    // 元のテキストでの終了位置
	TokenCount int               `json:"tokenCount"` // トークン数（概算）
	Index      int               `json:"index"`      // チャンクのインデックス
	Metadata   map[string]string `json:"metadata"`   // メタデータ情報
	Precedence int               `json:"precedence"` // チャンクの優先度
	References []string          `json:"references"` // 関連する見出しや参照情報
}

// DocumentChunker は拡張されたチャンキング実装を提供する
type DocumentChunker struct {
	config *ChunkConfig
}

func NewDocumentChunker(config *ChunkConfig) *DocumentChunker {
	if config == nil {
		config = NewDefaultConfig()
	}
	return &DocumentChunker{
		config: config,
	}
}

// ChunkDocument はドキュメントをチャンクに分割する
func (dc *DocumentChunker) ChunkDocument(content string) ([]Chunk, error) {
	// Front Matter の抽出と解析
	metadata, mainContent := dc.extractFrontMatter(content)

	// 現在の文字位置を追跡
	var currentPosition int

	// Front Matterのサイズを加算
	if len(metadata) > 0 {
		frontMatterEnd := strings.Index(content, "\n---\n")
		if frontMatterEnd != -1 {
			currentPosition = frontMatterEnd + 5 // \n---\n の長さを加算
		}
	}

	// セクション階層の解析
	sections := dc.parseSections(mainContent)

	var chunks []Chunk
	var headingStack []string // 見出し階層を追跡

	for _, section := range sections {
		// 見出しレベルに基づいて階層を更新
		headingStack = dc.updateHeadingStack(headingStack, section.level, section.title)

		// セクションの開始位置を記録
		sectionStartPos := currentPosition

		// セクションの分割
		sectionChunks := dc.chunkSection(section.content, dc.config.MaxTokens)

		// 各チャンクにメタデータと参照情報を付与
		for i, chunk := range sectionChunks {
			adjustedChunk := chunk
			adjustedChunk.Metadata = metadata
			adjustedChunk.References = make([]string, len(headingStack))
			copy(adjustedChunk.References, headingStack)

			// 正しい文字位置を設定
			if i == 0 {
				adjustedChunk.StartChar = sectionStartPos
			} else if len(chunks) > 0 {
				adjustedChunk.StartChar = chunks[len(chunks)-1].EndChar + 1
			}
			adjustedChunk.EndChar = adjustedChunk.StartChar + len(adjustedChunk.Content)

			// 優先度の計算
			adjustedChunk.Precedence = dc.calculatePrecedence(section.level, len(adjustedChunk.Content))

			adjustedChunk.Index = len(chunks)
			chunks = append(chunks, adjustedChunk)
		}

		// 現在位置を更新
		if len(chunks) > 0 {
			currentPosition = chunks[len(chunks)-1].EndChar + 1
		}

		// セクション間の区切りの長さを加算
		currentPosition += len(dc.config.ParagraphSeparator)
	}

	// チャンク間の意味的な関連性に基づく結合
	chunks = dc.mergeRelatedChunks(chunks)

	return chunks, nil
}

// updateHeadingStack は見出し階層を更新する
func (dc *DocumentChunker) updateHeadingStack(stack []string, level int, title string) []string {
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

// chunkSection はセクションの内容をチャンクに分割する
func (dc *DocumentChunker) chunkSection(content string, maxTokens int) []Chunk {
	content = strings.TrimSpace(content)
	if content == "" {
		return []Chunk{}
	}

	var chunks []Chunk
	paragraphs := strings.Split(content, dc.config.ParagraphSeparator)
	var currentChunk strings.Builder
	var currentTokens int

	// チャンクの追加用ヘルパー関数
	addChunk := func(content string, tokens int) {
		if content != "" {
			chunk := Chunk{
				Content:    strings.TrimSpace(content),
				TokenCount: tokens,
			}
			chunks = append(chunks, chunk)
		}
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		paraTokens := dc.countTokensJa(para)

		// 特殊な文書要素の重み付けを適用
		weight := 1.0
		if strings.HasPrefix(para, "- ") || strings.HasPrefix(para, "* ") {
			weight = dc.config.ListItemWeight
		} else if strings.HasPrefix(para, "```") {
			weight = dc.config.CodeBlockWeight
		} else if strings.Contains(para, "|") && strings.Contains(para, "-+-") {
			weight = dc.config.TableWeight
		}

		adjustedTokens := int(float64(paraTokens) * weight)

		// 見出し行の場合は新しいチャンクを開始
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

		// 現在のチャンクが最大トークン数を超える場合
		if currentTokens+adjustedTokens > maxTokens {
			if currentChunk.Len() > 0 {
				addChunk(currentChunk.String(), currentTokens)
				currentChunk.Reset()
				currentTokens = 0
			}
		}

		// 段落を追加
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentTokens += adjustedTokens

		// 日本語の文末表現でチャンクを分割するかを判断
		if dc.config.JapaneseConfig != nil {
			for _, ending := range dc.config.JapaneseConfig.SentenceEndings {
				if strings.HasSuffix(para, ending) && currentTokens >= dc.config.MinTokens {
					addChunk(currentChunk.String(), currentTokens)
					currentChunk.Reset()
					currentTokens = 0
					break
				}
			}
		}
	}

	// 残りのコンテンツをチャンクとして追加
	if currentChunk.Len() > 0 {
		addChunk(currentChunk.String(), currentTokens)
	}

	return chunks
}

// countTokensJa は日本語テキストのトークン数を計算する
func (dc *DocumentChunker) countTokensJa(text string) int {
	// 1. 全角スペースを半角スペースに変換
	text = strings.ReplaceAll(text, "　", " ")

	// 2. 記号による分割（句読点、括弧なども考慮）
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

	// 3. 空白で分割してカウント
	words := strings.Fields(text)

	// 4. 漢字・かな・カナの文字列をさらに分割してカウント
	total := 0
	for _, word := range words {
		if containsJapaneseCharacters(word) {
			// 日本語文字列は文字単位でカウント
			total += utf8.RuneCountInString(word)

			// 助詞やトピックマーカーの重み付け
			if dc.config.JapaneseConfig != nil {
				if isParticle(word) {
					total = int(float64(total) * dc.config.JapaneseConfig.KeyParticleWeight)
				}
				if isTopicMarker(word) {
					total = int(float64(total) * dc.config.JapaneseConfig.TopicMarkerWeight)
				}
			}
		} else {
			// 英数字などはワード単位でカウント
			total++
		}
	}

	return total
}

// 助詞かどうかを判定
func isParticle(word string) bool {
	particles := []string{"は", "が", "を", "に", "へ", "と", "で", "から", "まで", "より"}
	for _, p := range particles {
		if word == p {
			return true
		}
	}
	return false
}

// トピックマーカーかどうかを判定
func isTopicMarker(word string) bool {
	markers := []string{"は", "が", "について", "に関して"}
	for _, m := range markers {
		if word == m {
			return true
		}
	}
	return false
}

// 日本語文字かどうかを判定
func containsJapaneseCharacters(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana) ||
			unicode.In(r, unicode.Katakana) ||
			unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}

// extractFrontMatter は YAML Front Matter を抽出し解析する
func (dc *DocumentChunker) extractFrontMatter(content string) (map[string]string, string) {
	if !strings.HasPrefix(content, "---\n") {
		return make(map[string]string), content
	}

	// 2つ目の区切り文字を探す
	endIndex := strings.Index(content[4:], "\n---\n")
	if endIndex == -1 {
		return make(map[string]string), content
	}

	frontMatter := content[4 : endIndex+4] // 最初の"---\n"の後から次の"---"まで
	mainContent := content[endIndex+9:]    // 2つ目の"---\n"の後から

	// Front Matterのパース
	metadata := make(map[string]string)
	for _, line := range strings.Split(frontMatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(strings.TrimSpace(parts[0]), "\"")
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

		// 特別な処理が必要なフィールド
		if key == "tags" {
			// 配列形式の値を処理
			value = strings.Trim(value, "[]")
			metadata[key] = value
		} else {
			metadata[key] = value
		}
	}

	return metadata, mainContent
}

// セクション情報を保持する構造体
type section struct {
	title   string
	level   int
	content string
}

// parseSections はマークダウンからセクション階層を解析する
func (dc *DocumentChunker) parseSections(content string) []section {
	var sections []section
	lines := strings.Split(content, "\n")
	var currentSection section
	var buffer strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			// 前のセクションを保存
			if buffer.Len() > 0 {
				currentSection.content = buffer.String()
				sections = append(sections, currentSection)
				buffer.Reset()
			}

			// 新しいセクションの開始
			level := strings.Count(line, "#")
			title := strings.TrimSpace(strings.TrimLeft(line, "# "))
			currentSection = section{title: title, level: level}
		}
		buffer.WriteString(line + "\n")
	}

	// 最後のセクションを保存
	if buffer.Len() > 0 {
		currentSection.content = buffer.String()
		sections = append(sections, currentSection)
	}

	return sections
}

// calculatePrecedence はセクションの重要度を計算する
func (dc *DocumentChunker) calculatePrecedence(level int, contentLength int) int {
	// レベルが低いほど（より上位の見出しほど）優先度は高く
	// コンテンツが長いほど優先度は高く設定
	basePrecedence := 100 - (level * 20)
	contentFactor := contentLength / 100
	return basePrecedence + contentFactor
}

// mergeRelatedChunks は関連性の高いチャンクを結合する
func (dc *DocumentChunker) mergeRelatedChunks(chunks []Chunk) []Chunk {
	var merged []Chunk
	var current *Chunk

	for i := 0; i < len(chunks); i++ {
		if current == nil {
			current = &chunks[i]
			continue
		}

		// 関連性の判定
		if dc.areChunksRelated(*current, chunks[i]) &&
			current.TokenCount+chunks[i].TokenCount <= dc.config.MaxTokens {
			// チャンクの結合
			current.Content += "\n\n" + chunks[i].Content
			current.EndChar = chunks[i].EndChar
			current.TokenCount += chunks[i].TokenCount
			current.References = dc.mergeReferences(current.References, chunks[i].References)
		} else {
			merged = append(merged, *current)
			current = &chunks[i]
		}
	}

	if current != nil {
		merged = append(merged, *current)
	}

	return merged
}

// areChunksRelated は2つのチャンクの関連性を判定する
func (dc *DocumentChunker) areChunksRelated(chunk1, chunk2 Chunk) bool {
	// 同じ見出し階層に属しているか
	if len(chunk1.References) > 0 && len(chunk2.References) > 0 {
		return chunk1.References[0] == chunk2.References[0]
	}
	return false
}

// mergeReferences は参照情報をマージする
func (dc *DocumentChunker) mergeReferences(refs1, refs2 []string) []string {
	seen := make(map[string]bool)
	var merged []string

	for _, ref := range refs1 {
		if !seen[ref] {
			merged = append(merged, ref)
			seen[ref] = true
		}
	}

	for _, ref := range refs2 {
		if !seen[ref] {
			merged = append(merged, ref)
			seen[ref] = true
		}
	}

	return merged
}

// pkg/chunking/processor/processor.go
package processor

import (
	"fmt"
	"strings"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
)

// Processor はテキスト処理のインターフェースを定義
type Processor interface {
	Process(content string) ([]models.Chunk, error)
}

// DocumentProcessor はドキュメント処理の実装
type DocumentProcessor struct {
	config    *config.ChunkConfig
	extractor *FrontMatterExtractor
	chunker   *ContentChunker
	merger    *ChunkMerger
}

// NewDocumentProcessor は新しいDocumentProcessorを作成
func NewDocumentProcessor(cfg *config.ChunkConfig) *DocumentProcessor {
	return &DocumentProcessor{
		config:    cfg,
		extractor: NewFrontMatterExtractor(),
		chunker:   NewContentChunker(cfg),
		merger:    NewChunkMerger(cfg),
	}
}

// Process はドキュメントの処理を実行
func (p *DocumentProcessor) Process(content string) ([]models.Chunk, error) {
	// コンテンツの基本的なバリデーション
	if len(strings.TrimSpace(content)) == 0 {
		return nil, fmt.Errorf("empty content provided")
	}

	// Front Matterの抽出
	metadata, mainContent, err := p.extractor.Extract(content)
	if err != nil {
		return nil, fmt.Errorf("front matter extraction failed: %w", err)
	}

	// コンテンツのチャンキング
	chunks, err := p.chunker.Chunk(mainContent)
	if err != nil {
		return nil, fmt.Errorf("content chunking failed: %w", err)
	}

	// チャンクが生成されなかった場合のチェック
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks generated from content")
	}

	// メタデータとインデックスの設定
	currentPosition := 0
	for i := range chunks {
		// メタデータの設定
		for k, v := range metadata {
			chunks[i].SetMetadata(k, v)
		}

		// インデックスと位置情報の設定
		chunks[i].Index = i
		chunks[i].StartChar = currentPosition
		chunks[i].EndChar = currentPosition + len(chunks[i].Content)
		currentPosition = chunks[i].EndChar + len(p.config.ParagraphSeparator)
	}

	// チャンクの結合
	if !p.config.PreserveSections {
		chunks = p.merger.Merge(chunks)
	}

	return chunks, nil
}

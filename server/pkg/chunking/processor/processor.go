// pkg/chunking/processor/processor.go
package processor

import (
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
	// Front Matterの抽出
	metadata, mainContent, err := p.extractor.Extract(content)
	if err != nil {
		return nil, err
	}

	// コンテンツのチャンキング
	chunks, err := p.chunker.Chunk(mainContent)
	if err != nil {
		return nil, err
	}

	// メタデータの付与
	for i := range chunks {
		for k, v := range metadata {
			chunks[i].SetMetadata(k, v)
		}
	}

	// チャンクの結合
	mergedChunks := p.merger.Merge(chunks)

	return mergedChunks, nil
}

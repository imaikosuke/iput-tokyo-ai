// pkg/chunking/chunker.go
package chunking

import (
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/processor"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/utils"
)

// Chunker はパッケージの主要なインターフェース
type Chunker interface {
	ChunkDocument(content string) ([]models.Chunk, error)
	Configure(*config.ChunkConfig) error
	GetConfig() *config.ChunkConfig
}

// DocumentChunker は Chunker インターフェースの実装
type DocumentChunker struct {
	config    *config.ChunkConfig
	processor processor.Processor
}

// NewChunker は新しい Chunker インスタンスを作成
func NewChunker(cfg *config.ChunkConfig) (Chunker, error) {
	if cfg == nil {
		cfg = config.NewDefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		return nil, utils.NewChunkingError("initialization", "invalid configuration", err)
	}

	return &DocumentChunker{
		config:    cfg,
		processor: processor.NewDocumentProcessor(cfg),
	}, nil
}

// ChunkDocument はドキュメントをチャンクに分割
func (dc *DocumentChunker) ChunkDocument(content string) ([]models.Chunk, error) {
	if content == "" {
		return nil, utils.NewChunkingError("processing", "empty content", nil)
	}

	chunks, err := dc.processor.Process(content)
	if err != nil {
		return nil, utils.NewChunkingError("processing", "failed to process document", err)
	}

	return chunks, nil
}

// Configure は設定を更新
func (dc *DocumentChunker) Configure(cfg *config.ChunkConfig) error {
	if cfg == nil {
		return utils.NewChunkingError("configuration", "nil config provided", nil)
	}

	if err := cfg.Validate(); err != nil {
		return utils.NewChunkingError("configuration", "invalid configuration", err)
	}

	dc.config = cfg
	return nil
}

// GetConfig は現在の設定を取得
func (dc *DocumentChunker) GetConfig() *config.ChunkConfig {
	return dc.config
}

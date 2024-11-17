// pkg/chunking/processor/merger.go
package processor

import (
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/config"
	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/models"
)

// ChunkMerger はチャンクの結合を担当
type ChunkMerger struct {
	config *config.ChunkConfig
}

// NewChunkMerger は新しいChunkMergerを作成
func NewChunkMerger(cfg *config.ChunkConfig) *ChunkMerger {
	return &ChunkMerger{
		config: cfg,
	}
}

// Merge は関連するチャンクを結合
func (m *ChunkMerger) Merge(chunks []models.Chunk) []models.Chunk {
	if len(chunks) <= 1 {
		return chunks
	}

	var merged []models.Chunk
	var current *models.Chunk

	for i := 0; i < len(chunks); i++ {
		if current == nil {
			current = &chunks[i]
			continue
		}

		if m.canMerge(*current, chunks[i]) {
			// チャンクの結合
			current = m.mergeChunks(current, &chunks[i])
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

// canMerge は2つのチャンクが結合可能かどうかを判定
func (m *ChunkMerger) canMerge(chunk1, chunk2 models.Chunk) bool {
	// トークン数の合計が上限を超えないことを確認
	if chunk1.TokenCount+chunk2.TokenCount > m.config.MaxMergedTokens {
		return false
	}

	// 同じ見出し階層に属しているかを確認
	if len(chunk1.References) > 0 && len(chunk2.References) > 0 {
		return chunk1.References[0] == chunk2.References[0]
	}

	return false
}

// mergeChunks は2つのチャンクを結合する
func (m *ChunkMerger) mergeChunks(chunk1, chunk2 *models.Chunk) *models.Chunk {
	merged := &models.Chunk{
		Content:    chunk1.Content + "\n\n" + chunk2.Content,
		StartChar:  chunk1.StartChar,
		EndChar:    chunk2.EndChar,
		TokenCount: chunk1.TokenCount + chunk2.TokenCount,
		Index:      chunk1.Index,
		Metadata:   make(map[string]string),
		Precedence: (chunk1.Precedence + chunk2.Precedence) / 2,
	}

	// メタデータのマージ
	for k, v := range chunk1.Metadata {
		merged.Metadata[k] = v
	}
	for k, v := range chunk2.Metadata {
		if _, exists := merged.Metadata[k]; !exists {
			merged.Metadata[k] = v
		}
	}

	// 参照情報のマージ
	seenRefs := make(map[string]bool)
	for _, ref := range chunk1.References {
		if !seenRefs[ref] {
			merged.References = append(merged.References, ref)
			seenRefs[ref] = true
		}
	}
	for _, ref := range chunk2.References {
		if !seenRefs[ref] {
			merged.References = append(merged.References, ref)
			seenRefs[ref] = true
		}
	}

	return merged
}

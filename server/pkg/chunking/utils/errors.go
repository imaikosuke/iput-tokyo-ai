// pkg/chunking/utils/errors.go
package utils

import (
	"fmt"
)

// ChunkingError はチャンキング処理中のエラーを表す
type ChunkingError struct {
	Operation string
	Message   string
	Err       error
}

// Error はエラーメッセージを返す
func (e *ChunkingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Operation, e.Message)
}

// NewChunkingError は新しいChunkingErrorを作成
func NewChunkingError(operation, message string, err error) *ChunkingError {
	return &ChunkingError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

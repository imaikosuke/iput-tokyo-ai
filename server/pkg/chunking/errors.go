// pkg/chunking/errors.go
package chunking

import (
	"fmt"

	"github.com/imaikosuke/iput-tokyo-ai/server/pkg/chunking/utils"
)

// Error wraps all errors from the chunking package
type Error struct {
	*utils.ChunkingError
}

// WrapError wraps an error with additional context
func WrapError(err error, operation, message string) error {
	if err == nil {
		return nil
	}

	if chunkErr, ok := err.(*utils.ChunkingError); ok {
		return &Error{chunkErr}
	}

	return &Error{
		utils.NewChunkingError(operation, message, err),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("chunking error: %s", e.ChunkingError.Error())
}

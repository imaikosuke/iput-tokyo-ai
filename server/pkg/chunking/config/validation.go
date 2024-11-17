// pkg/chunking/config/validation.go
package config

import (
	"errors"
	"fmt"
)

// Validate は設定の妥当性を検証する
func (c *ChunkConfig) Validate() error {
	if c.MaxTokens < c.MinTokens {
		return errors.New("maxTokens must be greater than or equal to minTokens")
	}

	if c.OverlapTokens >= c.MaxTokens {
		return errors.New("overlapTokens must be less than maxTokens")
	}

	if c.MaxMergedTokens < c.MaxTokens {
		return errors.New("maxMergedTokens must be greater than or equal to maxTokens")
	}

	if err := c.validateWeights(); err != nil {
		return err
	}

	return nil
}

// validateWeights は重み付け設定の妥当性を検証する
func (c *ChunkConfig) validateWeights() error {
	weights := []float64{
		c.ListItemWeight,
		c.CodeBlockWeight,
		c.TableWeight,
		c.ContentFactor,
	}

	for i, w := range weights {
		if w <= 0 {
			return fmt.Errorf("weight at index %d must be positive", i)
		}
	}

	return nil
}

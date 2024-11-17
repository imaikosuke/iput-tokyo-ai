// pkg/chunking/processor/extractor.go
package processor

import (
	"strings"
)

// FrontMatterExtractor はFront Matterの抽出を担当
type FrontMatterExtractor struct{}

// NewFrontMatterExtractor は新しいFrontMatterExtractorを作成
func NewFrontMatterExtractor() *FrontMatterExtractor {
	return &FrontMatterExtractor{}
}

// Extract はFront Matterとメインコンテンツを抽出
func (e *FrontMatterExtractor) Extract(content string) (map[string]string, string, error) {
	if !strings.HasPrefix(content, "---\n") {
		return make(map[string]string), content, nil
	}

	endIndex := strings.Index(content[4:], "\n---\n")
	if endIndex == -1 {
		return make(map[string]string), content, nil
	}

	frontMatter := content[4 : endIndex+4]
	mainContent := content[endIndex+9:]

	metadata := e.parseFrontMatter(frontMatter)
	return metadata, mainContent, nil
}

// parseFrontMatter はFront Matterをパース
func (e *FrontMatterExtractor) parseFrontMatter(content string) map[string]string {
	metadata := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		metadata[key] = strings.Trim(value, "\"")
	}

	return metadata
}

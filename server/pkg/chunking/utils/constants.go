// pkg/chunking/utils/constants.go
package utils

const (
	// 最小チャンクサイズ
	MinChunkSize = 100

	// デフォルトの区切り文字
	DefaultSeparator = "\n\n"

	// 見出しの最大レベル
	MaxHeadingLevel = 6

	// 文字種別の重み付け
	JapaneseCharacterWeight = 1.0
	LatinCharacterWeight    = 0.5
	SymbolWeight            = 0.3
)

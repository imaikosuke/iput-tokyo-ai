// pkg/chunking/models/chunk.go
package models

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

// NewChunk は新しいChunkインスタンスを作成する
func NewChunk(content string) *Chunk {
    return &Chunk{
        Content:    content,
        Metadata:   make(map[string]string),
        References: make([]string, 0),
    }
}

// SetPosition はチャンクの開始位置と終了位置を設定する
func (c *Chunk) SetPosition(start, end int) {
    c.StartChar = start
    c.EndChar = end
}

// AddReference は参照情報を追加する
func (c *Chunk) AddReference(ref string) {
    c.References = append(c.References, ref)
}

// SetMetadata はメタデータを設定する
func (c *Chunk) SetMetadata(key, value string) {
    c.Metadata[key] = value
}

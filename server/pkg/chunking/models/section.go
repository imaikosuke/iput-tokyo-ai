// pkg/chunking/models/section.go
package models

// Section はドキュメントのセクション情報を表す構造体
type Section struct {
	Title   string // セクションのタイトル
	Level   int    // 見出しレベル
	Content string // セクションの内容
}

// NewSection は新しいSectionインスタンスを作成する
func NewSection(title string, level int, content string) *Section {
	return &Section{
		Title:   title,
		Level:   level,
		Content: content,
	}
}

// IsEmpty はセクションが空かどうかを判定する
func (s *Section) IsEmpty() bool {
	return s.Content == ""
}

// ContainsSubsection はサブセクションを含むかどうかを判定する
func (s *Section) ContainsSubsection() bool {
	return len(s.Content) > 0 && s.Content[0] == '#'
}

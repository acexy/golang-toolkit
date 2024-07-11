package str

import (
	"strings"
	"unicode"
)

// CharLength 返回字符串的字符长度
func CharLength(str string) int {
	return len([]rune(str))
}

// Substring 字符串截取指定长度 包含start下标，不含end下标
func Substring(str string, start int, end ...int) string {
	raw := []rune(str)
	var e int
	if len(end) > 0 {
		e = end[0]
	} else {
		e = CharLength(str)
	}
	return string(raw[start:e])
}

// HasText 检查是否具有有意义的内容
func HasText(str string) bool {
	return CharLength(strings.TrimSpace(str)) > 0
}

// HasLength 检查是否没有任何长度的值
func HasLength(str string) bool {
	return str != ""
}

type Builder struct {
	b strings.Builder
}

// NewBuilder 创建一个StringBuilder
func NewBuilder(str ...string) *Builder {
	builder := &Builder{}
	if len(str) > 0 {
		builder.b.WriteString(str[0])
	}
	return builder
}

func (b *Builder) WriteString(str string) *Builder {
	b.b.WriteString(str)
	return b
}

func (b *Builder) WriteByte(byte byte) *Builder {
	_ = b.b.WriteByte(byte)
	return b
}

func (b *Builder) WriteBytes(bytes []byte) *Builder {
	b.b.Write(bytes)
	return b
}
func (b *Builder) WriteRune(r rune) *Builder {
	b.b.WriteRune(r)
	return b
}

func (b *Builder) ToString() string {
	return b.b.String()
}

// LowFirstChar 首字母小写
func LowFirstChar(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

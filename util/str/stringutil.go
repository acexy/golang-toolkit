package str

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/acexy/golang-toolkit/math/conversion"
	"github.com/iancoleman/strcase"
)

// CharLength counts Unicode runes.
func CharLength(str string) int {
	return len([]rune(str))
}

// Substring slices by rune indexes with an exclusive end; invalid bounds panic.
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

// HasText reports whether non-whitespace text remains.
func HasText(str string) bool {
	return CharLength(strings.TrimSpace(str)) > 0
}

// Builder builds strings with chainable typed writes.
type Builder struct {
	b strings.Builder
}

// NewBuilder creates a builder with an optional initial string.
func NewBuilder(str ...string) *Builder {
	builder := &Builder{}
	if len(str) > 0 {
		builder.b.WriteString(str[0])
	}
	return builder
}

// WriteString appends a string and returns the builder.
func (b *Builder) WriteString(str string) *Builder {
	b.b.WriteString(str)
	return b
}

// WriteBool appends a boolean and returns the builder.
func (b *Builder) WriteBool(bool bool) *Builder {
	b.b.WriteString(strconv.FormatBool(bool))
	return b
}

// WriteInt8 appends a int8 as decimal text and returns the builder.
func (b *Builder) WriteInt8(i8 int8) *Builder {
	b.b.WriteString(conversion.FromInt8(i8))
	return b
}

// WriteUint8 appends a uint8 as decimal text and returns the builder.
func (b *Builder) WriteUint8(u8 uint8) *Builder {
	b.b.WriteString(conversion.FromUint8(u8))
	return b
}

// WriteInt16 appends a int16 as decimal text and returns the builder.
func (b *Builder) WriteInt16(i16 int16) *Builder {
	b.b.WriteString(conversion.FromInt16(i16))
	return b
}

// WriteUint16 appends a uint16 as decimal text and returns the builder.
func (b *Builder) WriteUint16(u16 uint16) *Builder {
	b.b.WriteString(conversion.FromUint16(u16))
	return b
}

// WriteInt32 appends a int32 as decimal text and returns the builder.
func (b *Builder) WriteInt32(i32 int32) *Builder {
	b.b.WriteString(conversion.FromInt32(i32))
	return b
}

// WriteUint32 appends a uint32 as decimal text and returns the builder.
func (b *Builder) WriteUint32(u32 uint32) *Builder {
	b.b.WriteString(conversion.FromUint32(u32))
	return b
}

// WriteInt appends a int as decimal text and returns the builder.
func (b *Builder) WriteInt(i int) *Builder {
	b.b.WriteString(conversion.FromInt(i))
	return b
}

// WriteUint appends a uint as decimal text and returns the builder.
func (b *Builder) WriteUint(u uint) *Builder {
	b.b.WriteString(conversion.FromUint(u))
	return b
}

// WriteInt64 appends a int64 as decimal text and returns the builder.
func (b *Builder) WriteInt64(i64 int64) *Builder {
	b.b.WriteString(conversion.FromInt64(i64))
	return b
}

// WriteUint64 appends a uint64 as decimal text and returns the builder.
func (b *Builder) WriteUint64(u64 uint64) *Builder {
	b.b.WriteString(conversion.FromUint64(u64))
	return b
}

// WriteFlot32 appends a float32 as decimal text and returns the builder.
func (b *Builder) WriteFlot32(f32 float32) *Builder {
	b.b.WriteString(conversion.FromFloat32(f32))
	return b
}

// WriteFlot64 appends a float64 as decimal text and returns the builder.
func (b *Builder) WriteFlot64(f64 float64) *Builder {
	b.b.WriteString(conversion.FromFloat64(f64))
	return b
}

// WriteByte appends one byte and returns the builder.
func (b *Builder) WriteByte(byte byte) *Builder {
	_ = b.b.WriteByte(byte)
	return b
}

// WriteBytes appends bytes and returns the builder.
func (b *Builder) WriteBytes(bytes []byte) *Builder {
	b.b.Write(bytes)
	return b
}

// WriteRune appends one Unicode rune and returns the builder.
func (b *Builder) WriteRune(r rune) *Builder {
	b.b.WriteRune(r)
	return b
}

// ToString returns the accumulated string.
func (b *Builder) ToString() string {
	return b.b.String()
}

// LowFirstChar lowercases the first Unicode rune.
func LowFirstChar(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// UpFirstChar uppercases the first Unicode rune.
func UpFirstChar(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// IsCamelCase checks the restricted ASCII lower-camel naming pattern.
func IsCamelCase(s string) bool {
	camelCasePattern := `^[a-z]+(?:[A-Z][a-z0-9]*)*$`
	re := regexp.MustCompile(camelCasePattern)
	return re.MatchString(s)
}

// IsSnakeCase checks the restricted ASCII lowercase underscore pattern.
func IsSnakeCase(s string) bool {
	re := regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)
	return re.MatchString(s)
}

// CamelToSnake converts a name to snake case using strcase.
func CamelToSnake(s string) string {
	return strcase.ToSnake(s)
	//re := regexp.MustCompile(`([A-Z])`)
	//snake := re.ReplaceAllString(s, `_$1`)
	//return strings.ToLower(snake)
}

// SnakeToCamel converts a name to lower camel case using strcase.
func SnakeToCamel(s string) string {
	return strcase.ToLowerCamel(s)
	//parts := strings.Split(s, "_")
	//for i := 0; i < len(parts); i++ {
	//	parts[i] = UpFirstChar(parts[i])
	//}
	//return LowFirstChar(strings.Join(parts, ""))
}

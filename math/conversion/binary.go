package conversion

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Binary 定义二进制值
type Binary struct {
	value int64
}

// Binaries 定义二进制数组
type Binaries []*Binary

// NewFromRawBinary 通过原始二进制值转换
// 010010110101
func NewFromRawBinary(rawBinary string) (*Binary, error) {
	decimal, err := strconv.ParseInt(rawBinary, 2, 64)
	if err != nil {
		return nil, err
	}
	return &Binary{value: decimal}, nil
}

// NewFromRawHex 通过指定的原始hex值转换
// ff
func NewFromRawHex(rawHex string) (*Binary, error) {
	rawHex = strings.ToLower(rawHex)
	decimal, err := strconv.ParseInt(strings.Replace(rawHex, hexPrefix, "", 1), 16, 64)
	if err != nil {
		return nil, err
	}
	return &Binary{value: decimal}, nil
}

// NewFromDecimal 通过指定十进制数据转换
func NewFromDecimal(decimal int64) *Binary {
	return &Binary{value: decimal}
}

// NewFromByte 通过byte转换
func NewFromByte(decimal byte) *Binary {
	return &Binary{value: int64(decimal)}
}

// Value 输出二进制值
func (b *Binary) Value() string {
	return fmt.Sprintf("%b", b.value)
}

// To8Bit 将该二进制值转换为计算机bit 8位
func (b *Binary) To8Bit() (string, error) {
	if b.value > 256 {
		return "", errors.New("binary value more than " + strconv.Itoa(maxBit))
	}
	return appendLeftZero(b.Value(), 8), nil
}

// ToHexValue 转换位hex值
func (b *Binary) ToHexValue() string {
	return fmt.Sprintf("%x", b.value)
}

// To2Hex 将该二进制值转化为hex 2位
func (b *Binary) To2Hex() (string, error) {
	if b.value > 256 {
		return "", errors.New("binary value more than " + strconv.Itoa(maxBit))
	}
	return appendLeftZero(b.ToHexValue(), 2), nil
}

// NewFormBytes 通过bytes转换
func NewFormBytes(bytes []byte) *Binaries {
	bs := make([]*Binary, len(bytes))
	for i, v := range bytes {
		bs[i] = NewFromByte(v)
	}
	bis := Binaries(bs)
	return &bis
}

func (b *Binaries) rawType() []*Binary {
	return *b
}

func (b *Binaries) join(valuer func(b *Binary) (string, error), char ...string) (string, error) {
	var builder strings.Builder
	bSlice := b.rawType()
	for i, v := range bSlice {
		value, err := valuer(v)
		if err != nil {
			return "", err
		}
		builder.WriteString(value)
		if len(char) > 0 {
			if i < len(bSlice)-1 {
				builder.WriteString(char[0])
			}
		}
	}
	return builder.String(), nil
}

// To8Bits 转换为二进制bit串 可以指定每8bit之间的连接填充字符
func (b *Binaries) To8Bits(char ...string) (string, error) {
	return b.join(func(b *Binary) (string, error) {
		value, err := b.To8Bit()
		if err != nil {
			return "error", err
		}
		return value, nil
	}, char...)
}

// To2HexString ToHexString 转换为hex字符串 可以指定每个hex值之间的连接填充字符
func (b *Binaries) To2HexString(char ...string) string {
	hexString, _ := b.join(func(b *Binary) (string, error) {
		value, err := b.To2Hex()
		if err != nil {
			return "error", err
		}
		return value, nil
	}, char...)
	return hexString
}

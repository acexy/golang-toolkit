package conversion

import (
	"fmt"
	"strconv"
	"strings"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

// Binary stores an int64 with binary and hexadecimal representations.
type Binary struct {
	value int64
}

// Binaries holds a sequence of binary values.
type Binaries []*Binary

// NewFromRawBinary parses a base-2 int64 string.
func NewFromRawBinary(rawBinary string) (*Binary, error) {
	decimal, err := strconv.ParseInt(rawBinary, 2, 64)
	if err != nil {
		return nil, err
	}
	return &Binary{value: decimal}, nil
}

// NewFromRawHex parses a hexadecimal int64 string with an optional 0x prefix.
func NewFromRawHex(rawHex string) (*Binary, error) {
	rawHex = strings.ToLower(rawHex)
	decimal, err := strconv.ParseInt(strings.Replace(rawHex, hexPrefix, "", 1), 16, 64)
	if err != nil {
		return nil, err
	}
	return &Binary{value: decimal}, nil
}

// NewFromDecimal creates a binary value from an int64.
func NewFromDecimal(decimal int64) *Binary {
	return &Binary{value: decimal}
}

// NewFromByte creates a binary value from a byte.
func NewFromByte(decimal byte) *Binary {
	return &Binary{value: int64(decimal)}
}

// Value returns the unpadded binary representation.
func (b *Binary) Value() string {
	return fmt.Sprintf("%b", b.value)
}

// To8Bit left-pads to eight characters and rejects values greater than 255.
func (b *Binary) To8Bit() (string, error) {
	if b.value > maxByte {
		return "", toolkitError.ErrBinaryValueOutOfByteRange
	}
	return appendLeftZero(b.Value(), 8), nil
}

// ToHexValue returns the unpadded hexadecimal representation.
func (b *Binary) ToHexValue() string {
	return fmt.Sprintf("%x", b.value)
}

// To2Hex left-pads to two hex characters and rejects values greater than 255.
func (b *Binary) To2Hex() (string, error) {
	if b.value > maxByte {
		return "", toolkitError.ErrBinaryValueOutOfByteRange
	}
	return appendLeftZero(b.ToHexValue(), 2), nil
}

// NewFromBytes creates one binary value for each byte.
func NewFromBytes(bytes []byte) *Binaries {
	bs := make([]*Binary, len(bytes))
	for i, v := range bytes {
		bs[i] = NewFromByte(v)
	}
	bis := Binaries(bs)
	return &bis
}

// RawType returns the underlying slice.
func (b *Binaries) RawType() []*Binary {
	return *b
}

func (b *Binaries) join(valuer func(b *Binary) (string, error), char ...string) (string, error) {
	var builder strings.Builder
	bSlice := b.RawType()
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

// To8Bits joins eight-character binary values with an optional separator.
func (b *Binaries) To8Bits(char ...string) (string, error) {
	return b.join(func(b *Binary) (string, error) {
		value, err := b.To8Bit()
		if err != nil {
			return "", err
		}
		return value, nil
	}, char...)
}

// To2Hex joins two-character hex values with an optional separator, discarding conversion errors.
func (b *Binaries) To2Hex(char ...string) string {
	hexString, _ := b.join(func(b *Binary) (string, error) {
		value, err := b.To2Hex()
		if err != nil {
			return "", err
		}
		return value, nil
	}, char...)
	return hexString
}

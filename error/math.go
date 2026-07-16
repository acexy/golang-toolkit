package error

import "errors"

var (
	// ErrBinaryValueOutOfByteRange 表示二进制值超出单字节范围
	ErrBinaryValueOutOfByteRange = errors.New("binary value out of byte range")

	// ErrEmptyProbability 表示概率参数为空
	ErrEmptyProbability = errors.New("empty probability")

	// ErrInvalidProbabilityTotal 表示概率总和无效
	ErrInvalidProbabilityTotal = errors.New("invalid probability total")

	// ErrProbabilityResultNotFound 表示未匹配到概率结果
	ErrProbabilityResultNotFound = errors.New("probability result not found")
)

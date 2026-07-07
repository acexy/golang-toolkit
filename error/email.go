package error

import "errors"

var (
	// ErrBadEmailClient 表示邮件客户端参数无效
	ErrBadEmailClient = errors.New("bad email client")

	// ErrEmptyToAddresses 表示收件人地址为空
	ErrEmptyToAddresses = errors.New("empty to addresses")

	// ErrBadEmailContent 表示邮件内容参数无效
	ErrBadEmailContent = errors.New("bad email content")

	// ErrCreateEmailClient 表示创建邮件客户端失败
	ErrCreateEmailClient = errors.New("failed to create email client")
)

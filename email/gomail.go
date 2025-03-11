package email

import (
	"errors"
	"gopkg.in/gomail.v2"
	"strings"
	"sync"
)

var goMailOnce sync.Once
var goMail *GoMail

type GoMail struct {
	dialer      *gomail.Dialer
	fromAddress string
}

type ToAddress struct {

	// 收件人地址 *
	Address string

	// 收件人称呼
	Name string
}

type Content struct {
	// 接收地址
	toAddresses []*ToAddress

	// 邮件标题
	subject string

	contentType string
	body        string

	// 附件文件路径
	attachments []string
}

func NewContent(toAddresses []*ToAddress, subject string) *Content {
	return &Content{
		toAddresses: toAddresses,
		subject:     subject,
	}
}

func (c *Content) SetContent(contentType, body string) *Content {
	c.contentType = contentType
	c.body = body
	return c
}

func (c *Content) SetAttach(attach []string) *Content {
	if len(attach) != 0 {
		c.attachments = attach
	}
	return c
}

func (c *Content) toMessage() (*gomail.Message, error) {
	if len(c.toAddresses) == 0 {
		return nil, errors.New("empty to addresses")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", goMail.fromAddress)
	if c.subject != "" {
		message.SetHeader("subject", c.subject)
	}
	addressWithoutMemo := make([]string, 0)
	addressWithMemo := make([]string, 0)
	for _, toAddress := range c.toAddresses {
		if toAddress.Name == "" {
			addressWithoutMemo = append(addressWithoutMemo, toAddress.Address)
		} else {
			addressWithMemo = append(addressWithMemo, toAddress.Name+":"+toAddress.Address)
		}
	}

	if len(addressWithoutMemo) > 0 {
		message.SetHeader("To", addressWithoutMemo...)
	}

	if len(addressWithMemo) > 0 {
		var setTo bool
		for _, addr := range addressWithMemo {
			split := strings.Split(addr, ":")
			if len(addressWithoutMemo) == 0 && !setTo {
				message.SetHeader("To", split[1])
				setTo = true
			}
			message.SetAddressHeader(split[0], split[1], split[0])
		}

	}

	message.SetBody(c.contentType, c.body)

	if len(c.attachments) > 0 {
		for _, attachment := range c.attachments {
			message.Attach(attachment)
		}
	}

	return message, nil
}

func NewGoMail(host string, port int, username, password string, fromAddress string, isSSL ...bool) *GoMail {
	goMailOnce.Do(func() {
		dialer := gomail.NewDialer(host, port, username, password)
		if len(isSSL) > 0 && isSSL[0] {
			dialer.SSL = true
		}
		goMail = &GoMail{
			dialer:      dialer,
			fromAddress: fromAddress,
		}
	})
	return goMail
}

func (g *GoMail) SendMail(content *Content) error {
	if content == nil {
		return errors.New("bad content param")
	}
	m, err := content.toMessage()
	if err != nil {
		return err
	}
	return g.dialer.DialAndSend(m)
}

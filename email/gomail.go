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

	return message, nil
}

func NewGoMail(host string, port int, username, password string, fromAddress string) *GoMail {
	goMailOnce.Do(func() {
		goMail = &GoMail{
			dialer:      gomail.NewDialer(host, port, username, password),
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

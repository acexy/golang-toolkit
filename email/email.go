package email

import (
	toolkitError "github.com/acexy/golang-toolkit/error"
	mail "github.com/wneessen/go-mail"
)

type Client struct {
	client          *mail.Client
	initErr         error
	fromEmail       string
	fromDisplayName string
}

type Address struct {
	// 邮件地址
	Email string
	// 显示名称
	DisplayName string
}

type Message struct {

	// 发送方地址 可选
	fromEmail string
	// 发送方名称 可选
	fromDisplayName string

	// 接收地址
	toAddresses []*Address
	// 邮件标题
	subject         string
	bodyContentType string
	body            string

	// 附件文件路径
	attachments []string
}

func NewMessage(toAddresses []*Address, subject string) *Message {
	return &Message{
		toAddresses: toAddresses,
		subject:     subject,
	}
}

func (m *Message) SetBody(contentType, body string) *Message {
	m.bodyContentType = contentType
	m.body = body
	return m
}

func (m *Message) SetFrom(fromEmail, fromDisplayName string) *Message {
	m.fromEmail = fromEmail
	m.fromDisplayName = fromDisplayName
	return m
}

func (m *Message) SetAttachments(attachments []string) *Message {
	if len(attachments) != 0 {
		m.attachments = attachments
	}
	return m
}

func (m *Message) toMessage(c *Client) (*mail.Msg, error) {
	if c == nil {
		return nil, toolkitError.ErrBadEmailClient
	}
	if len(m.toAddresses) == 0 {
		return nil, toolkitError.ErrEmptyToAddresses
	}
	message := mail.NewMsg()
	fromEmail := c.fromEmail
	fromDisplayName := c.fromDisplayName
	if m.fromEmail != "" {
		fromEmail = m.fromEmail
		fromDisplayName = m.fromDisplayName
	}
	if fromDisplayName != "" {
		if err := message.FromFormat(fromDisplayName, fromEmail); err != nil {
			return nil, err
		}
	} else if err := message.From(fromEmail); err != nil {
		return nil, err
	}
	for _, toAddress := range m.toAddresses {
		if toAddress == nil || toAddress.Email == "" {
			continue
		}
		if toAddress.DisplayName == "" {
			if err := message.AddTo(toAddress.Email); err != nil {
				return nil, err
			}
		} else {
			if err := message.AddToFormat(toAddress.DisplayName, toAddress.Email); err != nil {
				return nil, err
			}
		}
	}
	recipients, err := message.GetRecipients()
	if err != nil || len(recipients) == 0 {
		return nil, toolkitError.ErrEmptyToAddresses
	}
	if m.subject != "" {
		message.Subject(m.subject)
	}
	contentType := mail.TypeTextPlain
	if m.bodyContentType != "" {
		contentType = mail.ContentType(m.bodyContentType)
	}
	message.SetBodyString(contentType, m.body)
	if len(m.attachments) > 0 {
		for _, attachment := range m.attachments {
			message.AttachFile(attachment)
		}
	}
	return message, nil
}

func NewClient(host string, port int, username, password, fromEmail string, useSSL bool) *Client {
	return NewClientWithName(host, port, username, password, fromEmail, "", useSSL)
}

func NewClientWithName(host string, port int, username, password, fromEmail, fromDisplayName string, useSSL bool) *Client {
	options := []mail.Option{
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
	}
	if useSSL {
		options = append(options, mail.WithSSL())
	}
	if username != "" || password != "" {
		options = append(options, mail.WithSMTPAuth(mail.SMTPAuthPlain))
	}
	client, err := mail.NewClient(host, options...)
	return &Client{
		client:          client,
		initErr:         err,
		fromEmail:       fromEmail,
		fromDisplayName: fromDisplayName,
	}
}

func (c *Client) SendMail(message *Message) error {
	if message == nil {
		return toolkitError.ErrBadEmailContent
	}
	if c == nil || c.initErr != nil || c.client == nil {
		return toolkitError.ErrCreateEmailClient
	}
	m, err := message.toMessage(c)
	if err != nil {
		return err
	}
	return c.client.DialAndSend(m)
}

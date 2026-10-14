package email

import (
	toolkitError "github.com/acexy/golang-toolkit/error"
	mail "github.com/wneessen/go-mail"
)

// Client wraps an SMTP client and its default sender configuration.
type Client struct {
	client          *mail.Client
	initErr         error
	fromEmail       string
	fromDisplayName string
}

// Address pairs an email address with an optional display name.
type Address struct {
	// Email is the recipient email address.
	Email string
	// DisplayName is the optional recipient display name.
	DisplayName string
}

// Message holds recipients, subject, body, sender overrides, and attachment paths.
type Message struct {

	// Optional sender address override.
	fromEmail string
	// Optional sender display name override.
	fromDisplayName string

	// Recipient addresses.
	toAddresses []*Address
	// Message subject.
	subject         string
	bodyContentType string
	body            string

	// Attachment file paths.
	attachments []string
}

// NewMessage creates a message with the supplied recipients and subject.
func NewMessage(toAddresses []*Address, subject string) *Message {
	return &Message{
		toAddresses: toAddresses,
		subject:     subject,
	}
}

// SetBody sets the body and content type; an empty type uses plain text.
func (m *Message) SetBody(contentType, body string) *Message {
	m.bodyContentType = contentType
	m.body = body
	return m
}

// SetFrom overrides the client sender when fromEmail is nonempty.
func (m *Message) SetFrom(fromEmail, fromDisplayName string) *Message {
	m.fromEmail = fromEmail
	m.fromDisplayName = fromDisplayName
	return m
}

// SetAttachments sets attachment paths; an empty slice leaves existing attachments unchanged.
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

// NewClient creates an SMTP wrapper; initialization errors are reported by SendMail.
func NewClient(host string, port int, username, password, fromEmail string, useSSL bool) *Client {
	return NewClientWithName(host, port, username, password, fromEmail, "", useSSL)
}

// NewClientWithName creates an SMTP wrapper with a sender display name; errors are reported by SendMail.
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

// SendMail validates the message and client, then sends it through SMTP.
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

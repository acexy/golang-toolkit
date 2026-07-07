package email

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMessageToMailMessage(t *testing.T) {
	client := NewClientWithName("smtp.example.com", 465, "from@example.com", "password", "from@example.com", "from", true)

	message, err := NewMessage([]*Address{
		{Email: "to1@example.com", DisplayName: "to1"},
		{Email: "to2@example.com"},
	}, "test").SetBody("text/plain", "test").toMessage(client)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err = message.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	content := buf.String()
	if !strings.Contains(content, `From: "from" <from@example.com>`) {
		t.Fatalf("missing From header: %s", content)
	}
	if !strings.Contains(content, `To: "to1" <to1@example.com>, <to2@example.com>`) {
		t.Fatalf("missing To header: %s", content)
	}
	if !strings.Contains(content, "Subject: test") {
		t.Fatalf("missing Subject header: %s", content)
	}
}

func TestSendMail(t *testing.T) {
	if os.Getenv("GOLANG_TOOLKIT_EMAIL_SEND_TEST") != "true" {
		t.Skip("set GOLANG_TOOLKIT_EMAIL_SEND_TEST=true to run real email sending test")
	}

	port, err := strconv.Atoi(os.Getenv("GOLANG_TOOLKIT_EMAIL_PORT"))
	if err != nil {
		t.Fatal(err)
	}

	client := NewClientWithName(
		os.Getenv("GOLANG_TOOLKIT_EMAIL_HOST"),
		port,
		os.Getenv("GOLANG_TOOLKIT_EMAIL_USERNAME"),
		os.Getenv("GOLANG_TOOLKIT_EMAIL_PASSWORD"),
		os.Getenv("GOLANG_TOOLKIT_EMAIL_FROM"),
		os.Getenv("GOLANG_TOOLKIT_EMAIL_FROM_NAME"),
		os.Getenv("GOLANG_TOOLKIT_EMAIL_SSL") == "true",
	)

	err = client.SendMail(NewMessage([]*Address{
		{
			Email:       os.Getenv("GOLANG_TOOLKIT_EMAIL_TO"),
			DisplayName: os.Getenv("GOLANG_TOOLKIT_EMAIL_TO_NAME"),
		},
	}, "golang-toolkit email test").
		SetBody("text/plain", "golang-toolkit email test").
		SetAttachments([]string{"./email_test.go"}),
	)
	if err != nil {
		t.Fatal(err)
	}
}

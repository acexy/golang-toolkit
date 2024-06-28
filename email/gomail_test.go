package email

import "testing"

func TestSendMail(t *testing.T) {

	mail := NewGoMail("smtp.exmail.qq.com", 465, "from@example.com", "password", "from@example.com")

	err := mail.SendMail(NewContent([]*ToAddress{
		{Address: "to1@example.com", Name: "to1"},
		{Address: "to2@example.com"},
	}, "test").SetContent("text/plain", "test").SetAttach([]string{"/Users/acexy/Downloads/newlog"}))
	if err != nil {
		t.Error("send mail error", err)
	}
}

package helpers

import (
	"github.com/scorredoira/email"
	"net/mail"
	"strings"
)

func SendMail(subject *string, body *string, from *string, to *string, attachment *string, smtp *string) {
	payload := email.NewMessage(*subject, *body)
	payload.From = mail.Address{Name: "Backup", Address: *from}
	payload.To = strings.Split(*to, ",")
	if attachment != nil {
		err := payload.Attach(*attachment)
		LogError(err)
	}
	err := email.Send(*smtp, nil, payload)
	LogError(err)
}

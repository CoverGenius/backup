package notifier

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	time_format = "2006-01-02T15-04-05"
)

var (
	message_body_tmpl = `
{HEADER}

Job: {NAME}
Started: {START_TIME}
Finished: {END_TIME}
Duration: {DURATION}
`
)

type Mail struct {
	Subject *string
	Header  *string
	Body    *string
}

func (m *Mail) SetSubject(c *base.Config) error {
	var sbj string
	switch *c.Notifier.Status.Status {
	case "OK":
		sbj = fmt.Sprintf("%s %s", "[Backup::Success]", *c.Name)
	case "WARN":
		sbj = fmt.Sprintf("%s %s", "[Backup::Warning]", *c.Name)
	default:
		sbj = fmt.Sprintf("%s %s", "[Backup::Failure]", *c.Name)
	}
	m.Subject = &sbj

	return nil
}

func (m *Mail) SetHeader(c *base.Config) error {
	var hdr string
	switch *c.Notifier.Status.Status {
	case "OK":
		hdr = "Backup Completed Successfully!"
	case "WARN":
		hdr = "Backup Completed Successfully (with Warnings)!"
	default:
		hdr = "Backup Failed!"
	}
	m.Header = &hdr

	return nil
}

func (m *Mail) ComposeBody(c *base.Config) error {
	m.SetHeader(c)
	body := h.FormatString(
		&message_body_tmpl,
		"{HEADER}", *m.Header,
		"{NAME}", *c.Name,
		"{START_TIME}", *c.Notifier.Status.FormatStartTime(),
		"{END_TIME}", *c.Notifier.Status.FormatEndTime(),
		"{DURATION}", *c.Notifier.Status.DurationToString(),
	)
	m.Body = &body
	return nil
}

func (m *Mail) Verify(c *base.Config) error {
	if c.Notifier.Mail.To == nil || c.Notifier.Mail.From == nil {
		return errors.New("Either TO or FROM is not set for mail notifier!")
	}
	return nil
}

func (m *Mail) Pre(c *base.Config) error {
	return nil
}

func (m *Mail) Notify(c *base.Config) error {
	m.SetSubject(c)
	m.ComposeBody(c)

	if *c.Notifier.Status.Status == "FAILURE" && h.Log.Level != logrus.DebugLevel {
		body := fmt.Sprintf("%s\n\nSee the attached backup log for details.", *m.Body)
		h.SendMail(m.Subject, &body, c.Notifier.Mail.From, c.Notifier.Mail.To, c.LogFile, c.Notifier.Mail.Address)
	} else {
		h.SendMail(m.Subject, m.Body, c.Notifier.Mail.From, c.Notifier.Mail.To, nil, c.Notifier.Mail.Address)
	}
	return nil
}

func (m *Mail) Post(c *base.Config) error {
	return nil
}

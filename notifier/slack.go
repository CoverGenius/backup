package notifier

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Slack struct {
	Title   *string
	Color   *string
	Pretext *string
}

func (s *Slack) SetPretext(c *base.Config) error {
	var m string
	switch *c.Notifier.Status.Status {
	case "OK":
		m = fmt.Sprintf("%s %s", "[Backup::Success]", *c.Name)
	case "WARN":
		m = fmt.Sprintf("%s %s", "[Backup::Warning]", *c.Name)
	default:
		m = fmt.Sprintf("%s %s", "[Backup::Failure]", *c.Name)
	}
	s.Pretext = &m
	return nil
}

func (s *Slack) SetColor(c *base.Config) error {
	var m string
	switch *c.Notifier.Status.Status {
	case "OK":
		m = "good"
	case "WARN":
		m = "warning"
	default:
		m = "danger"
	}
	s.Color = &m
	return nil
}

func (s *Slack) SetTitle(c *base.Config) error {
	var m string
	switch *c.Notifier.Status.Status {
	case "OK":
		m = "Backup Completed Successfully!"
	case "WARN":
		m = "Backup Completed Successfully (with Warnings)!"
	default:
		m = "Backup Failed!"
	}
	s.Title = &m
	return nil
}

type Payload struct {
	Channel     *string      `json:"channel"`
	IconEmoji   *string      `json:"icon_emoji"`
	Username    *string      `json:"username"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text    *string `json:"text"`
	Color   *string `json:"color"`
	Pretext *string `json:"pretext"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Title *string `json:"title"`
	Value *string `json:"value"`
	Short *bool   `json:"short"`
}

func (s *Slack) Verify(c *base.Config) error {
	if c.Notifier.Slack.WebhookURL == nil || c.Notifier.Slack.Channel == nil {
		return errors.New("Either TO or FROM is not set for mail notifier!")
	}
	return nil
}

func (s *Slack) Pre(c *base.Config) error {
	if c.Notifier.Slack.Username == nil {
		c.Notifier.Slack.Username = h.StringP("Backup")
	}
	return nil
}

func ComposePayload(s *Slack, c *base.Config) *Payload {
	p := Payload{
		Channel:   c.Notifier.Slack.Channel,
		IconEmoji: h.StringP(":glitch_crab:"),
		Username:  c.Notifier.Slack.Username,
		Attachments: []Attachment{
			{
				Color:   s.Color,
				Pretext: s.Pretext,
				Text:    s.Title,
				Fields: []Field{
					{
						Title: h.StringP("Job"),
						Value: c.Name,
						Short: h.BoolP(false),
					},
					{
						Title: h.StringP("Started"),
						Value: c.Notifier.Status.FormatStartTime(),
						Short: h.BoolP(true),
					},
					{
						Title: h.StringP("Finished"),
						Value: c.Notifier.Status.FormatEndTime(),
						Short: h.BoolP(true),
					},
					{
						Title: h.StringP("Duration"),
						Value: c.Notifier.Status.DurationToString(),
						Short: h.BoolP(false),
					},
				},
			},
		},
	}
	if *c.Notifier.Status.Status == "FAILURE" && h.Log.Level != logrus.DebugLevel {
		log, _ := ioutil.ReadFile(*c.LogFile)
		log_s := fmt.Sprintf("%s", log)
		field := Field{
			Title: h.StringP("Detailed Backup Log"),
			Value: h.StringP(log_s),
			Short: h.BoolP(false),
		}
		p.Attachments[0].Fields = append(p.Attachments[0].Fields, field)
	}
	return &p
}

func (s *Slack) Notify(c *base.Config) error {
	s.SetPretext(c)
	s.SetColor(c)
	s.SetTitle(c)

	p := ComposePayload(s, c)
	var b bytes.Buffer
	h.JSONEncode(&b, p)

	h.MakeHTTPRequest(c.Notifier.Slack.WebhookURL, "POST", nil, nil, nil, &b, false)
	return nil
}

func (s *Slack) Post(c *base.Config) error {
	return nil
}

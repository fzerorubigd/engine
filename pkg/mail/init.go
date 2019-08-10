package mail

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"

	"elbix.dev/engine/pkg/assert"
	"elbix.dev/engine/pkg/config"
	"elbix.dev/engine/pkg/initializer"
	"elbix.dev/engine/pkg/job"
)

const mailerTopic = "mailer"

var (
	dialer *gomail.Dialer

	smtpUsername = config.RegisterString("mail.smtp.username", "", "smtp user name")
	smtpPassword = config.RegisterString("mail.smtp.password", "", "smtp password")

	smtpHost = config.RegisterString("mail.smtp.host", "0127.0.0.1", "smtp host")
	smtpPort = config.RegisterInt("mail.smtp.port", 1025, "smtp port")
)

// EmailAddress is the simple mail <mail@mail.com>
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Payload is the email payload for chapar
type Payload struct {
	Subject string         `json:"subject"`
	Message string         `json:"message"`
	From    EmailAddress   `json:"from"`
	To      []EmailAddress `json:"to"`
}

// NewEmailAddress return a new struct contain email address and name, just a shortcut
func NewEmailAddress(mail string) EmailAddress {
	return EmailAddress{
		Email: mail,
		Name:  "",
	}
}

// NewEmailNameAddress return a new struct contain email address and name, just a shortcut
func NewEmailNameAddress(mail, name string) EmailAddress {
	return EmailAddress{
		Email: mail,
		Name:  name,
	}
}

type setup struct {
}

func (setup) Process(ctx context.Context, data []byte) error {
	var msg Payload

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return errors.Wrap(err, "invalid payload")
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", msg.From.Email, msg.From.Name)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Message)

	toString := make([]string, len(msg.To))
	for i := range msg.To {
		toString[i] = m.FormatAddress(msg.To[i].Email, msg.To[i].Name)
	}
	m.SetHeader("To", toString...)

	return errors.Wrap(dialer.DialAndSend(m), "send failed")
}

// Send is actually enqueue them to send
func Send(ctx context.Context, subject, msg string, from EmailAddress, to ...EmailAddress) error {
	if len(to) < 1 {
		return errors.New("at least one recipient is required")
	}

	payload := Payload{
		Subject: subject,
		Message: msg,
		From:    from,
		To:      to,
	}

	data, err := json.Marshal(payload)
	assert.Nil(err)

	return job.EnqueueJob(ctx, mailerTopic, data)
}

func (s *setup) Initialize(context.Context) {
	dialer = gomail.NewDialer(
		smtpHost.String(), smtpPort.Int(), smtpUsername.String(), smtpPassword.String())

	job.RegisterWorker(mailerTopic, s)
}

func init() {
	initializer.Register(&setup{}, 0)
}

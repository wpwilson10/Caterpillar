package setup

import (
	"gopkg.in/gomail.v2"
)

// SendEmail sends a message to my email with the given subject and html body
func SendEmail(subject string, body string) {
	// setup email context
	m := gomail.NewMessage()
	m.SetHeader("From", "PatrickWilsonSR0@gmail.com")
	m.SetHeader("To", "wpwilson10@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// this is an app specific password
	d := gomail.NewDialer("smtp.gmail.com", 587, "PatrickWilsonSR0", "fjdczztkowcieddi")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		LogCommon(err).Error("Failed Dial and Send")
	}
}

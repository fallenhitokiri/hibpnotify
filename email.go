package hibpnotify

import (
	"log"
	"net/smtp"
)

type notificationEmail struct {
	config *config
}

func newNotificationEmail(c *config) (notification, error) {
	n := &notificationEmail{
		c,
	}
	return n, nil
}

func (n *notificationEmail) notify(accounts []string) error {
	if len(accounts) == 0 {
		log.Println("No breached accounts")
		return nil
	}

	return n.byEmail(accounts)
}

func (n *notificationEmail) byEmail(accounts []string) error {
	auth := smtp.PlainAuth("", n.config.NotifySMTPUser, n.config.NotifySMTPPassword, n.config.NotifySMTPHost)
	body, err := n.toEmailBody(accounts)

	if err != nil {
		return err
	}

	return smtp.SendMail(
		n.config.NotifySMTPAddr,
		auth,
		n.config.NotifySMTPFrom,
		[]string{n.config.NotifyEmail},
		[]byte(body),
	)
}

func (n *notificationEmail) toEmailBody(accounts []string) (string, error) {
	body := "Newly breached accounts:\n\n"

	for _, s := range accounts {
		body = body + s + "\n"
	}

	return body, nil
}

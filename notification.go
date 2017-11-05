package hibpnotify

type notification interface {
	notify(accounts []string) error
}

func newNotification(c *config) ([]notification, error) {
	n := []notification{}

	if c.NotifyEmail != "" {
		e, err := newNotificationEmail(c)

		if err != nil {
			return nil, err
		}

		n = append(n, e)
	}
	return n, nil
}

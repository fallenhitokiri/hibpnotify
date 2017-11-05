package hibpnotify

type notificationMock struct {
	reported []string
}

func (n *notificationMock) notify(accounts []string) error {
	n.reported = append(n.reported, accounts...)
	return nil
}

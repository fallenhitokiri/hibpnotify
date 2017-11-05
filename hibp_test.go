package hibpnotify

import (
	"time"
)

const accountBreachedOnce = "once@breached.tld"
const accountBreachedTwice = "twice@breached.tld"
const accountBreachedNever = "foo@bar.tld"

type hibpClientMock struct{}

func (hc *hibpClientMock) requestForEmail(email string) ([]*Breach, error) {
	breaches := []*Breach{}

	if email == accountBreachedOnce {
		b := &Breach{
			Title:        "title",
			Name:         "name",
			Domain:       "once.tld",
			AddedDate:    time.Now(),
			ModifiedDate: time.Now(),
			Description:  "abc",
		}
		breaches = append(breaches, b)
	} else if email == accountBreachedTwice {
		b := &Breach{
			Title:        "title",
			Name:         "name",
			Domain:       "once.tld",
			AddedDate:    time.Now(),
			ModifiedDate: time.Now(),
			Description:  "abc",
		}
		b2 := &Breach{
			Title:        "title2",
			Name:         "name2",
			Domain:       "twice.tld",
			AddedDate:    time.Now(),
			ModifiedDate: time.Now(),
			Description:  "abc",
		}
		breaches = append(breaches, b)
		breaches = append(breaches, b2)
	}

	return breaches, nil
}

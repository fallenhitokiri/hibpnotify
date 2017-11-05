package hibpnotify

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type updater struct {
	config       *config
	hibp         hibpClient
	db           *gorm.DB
	notification []notification
	dataSources  []dataSource
	newBreached  []string

	sleepTime     time.Duration
	frequencyTime time.Duration
}

func newUpdater(c *config, n []notification, ds []dataSource, hibp hibpClient) (*updater, error) {
	db, err := c.GetDatabase()

	if err != nil {
		return nil, err
	}

	up := &updater{
		config:        c,
		db:            db,
		hibp:          hibp,
		notification:  n,
		dataSources:   ds,
		sleepTime:     time.Duration(c.CheckRequestTimeout) * time.Second,
		frequencyTime: time.Duration(c.CheckIntervalInSeconds) * time.Second,
	}
	return up, nil
}

func (u *updater) update() {
	u.newBreached = []string{}

	for _, source := range u.dataSources {
		if err := u.handleSource(source); err != nil {
			log.Println(err)
		}
	}

	log.Println("newly breached: ", u.newBreached)
	for _, n := range u.notification {
		if err := n.notify(u.newBreached); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(u.frequencyTime)
}

func (u *updater) handleSource(source dataSource) error {
	ac, err := source.GetAccounts()

	if err != nil {
		return err
	}

	for _, email := range ac {
		log.Println("getting breachs for ", email)
		breaches, err := u.hibp.requestForEmail(email)

		if err != nil {
			return err
		}

		for _, breach := range breaches {
			b, err := u.saveBreach(breach)

			if err != nil {
				return err
			}

			if err := u.associateWithUser(b, email); err != nil {
				return err
			}
		}

		time.Sleep(u.sleepTime)
	}

	return nil
}

func (u *updater) saveBreach(breach *Breach) (*Breach, error) {
	var existing Breach
	err := u.db.Where("name = ?", breach.Name).First(&existing).Error

	if err == nil {
		return &existing, nil
	}

	err = u.db.Create(breach).Error
	return breach, err
}

func (u *updater) associateWithUser(breach *Breach, email string) error {
	var user User
	if err := u.db.Preload("Breaches").FirstOrCreate(&user, User{Email: email}).Error; err != nil {
		return err
	}

	for _, b := range user.Breaches {
		if b.Name == breach.Name && b.AddedDate == breach.AddedDate && b.Domain == breach.Domain &&
			b.Title == breach.Title {
			return nil // user already associated with breach
		}
	}

	u.userToNewBreached(email)
	return u.db.Model(&user).Association("Breaches").Append(breach).Error
}

func (u *updater) userToNewBreached(email string) {
	for _, ex := range u.newBreached {
		if ex == email {
			return
		}
	}

	u.newBreached = append(u.newBreached, email)
}

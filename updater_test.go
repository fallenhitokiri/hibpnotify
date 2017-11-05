package hibpnotify

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func TestUserToNewBreached(t *testing.T) {
	u := &updater{newBreached: []string{}}

	u.userToNewBreached("foo@bar.tld")

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	u.userToNewBreached("foo@bar.tld")

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	u.userToNewBreached("baz@bar.tld")

	if len(u.newBreached) != 2 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}
}

func TestAssociateWithUser(t *testing.T) {
	c := &config{
		DatabaseDialect: "sqlite3",
		DatabaseArgs:    ":memory:",
	}
	db, err := c.GetDatabase()

	if err != nil {
		t.Fatal(err)
	}

	u := &updater{
		db:          db,
		newBreached: []string{},
	}

	b := &Breach{
		Title:        "title",
		Name:         "name",
		Domain:       "once.tld",
		AddedDate:    time.Now(),
		ModifiedDate: time.Now(),
		Description:  "abc",
	}
	if err := db.Create(b).Error; err != nil {
		t.Fatal(err)
	}

	if err := u.associateWithUser(b, "foo@bar.tld"); err != nil {
		t.Fatal(err)
	}

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	if err := u.associateWithUser(b, "foo@bar.tld"); err != nil {
		t.Fatal(err)
	}

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	if err := u.associateWithUser(b, "baz@bar.tld"); err != nil {
		t.Fatal(err)
	}

	if len(u.newBreached) != 2 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}
}

func TestAssociateWithUserSecondBreach(t *testing.T) {
	c := &config{
		DatabaseDialect: "sqlite3",
		DatabaseArgs:    ":memory:",
	}
	db, err := c.GetDatabase()

	if err != nil {
		t.Fatal(err)
	}

	u := &updater{
		db:          db,
		newBreached: []string{},
	}

	b := &Breach{
		Title:        "title",
		Name:         "name",
		Domain:       "once.tld",
		AddedDate:    time.Now(),
		ModifiedDate: time.Now(),
		Description:  "abc",
	}
	if err := db.Create(b).Error; err != nil {
		t.Fatal(err)
	}

	if err := u.associateWithUser(b, "foo@bar.tld"); err != nil {
		t.Fatal(err)
	}

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	b2 := &Breach{
		Title:        "title2",
		Name:         "name2",
		Domain:       "twice.tld",
		AddedDate:    time.Now(),
		ModifiedDate: time.Now(),
		Description:  "abc",
	}
	if err := db.Create(b2).Error; err != nil {
		t.Fatal(err)
	}

	if err := u.associateWithUser(b2, "foo@bar.tld"); err != nil {
		t.Fatal(err)
	}

	if len(u.newBreached) != 1 {
		t.Fatal("newBreached wrong length ", len(u.newBreached), " ", u.newBreached)
	}

	var user User
	if err := db.Where("email = ?", "foo@bar.tld").Preload("Breaches").First(&user).Error; err != nil {
		t.Fatal(err)
	}

	if len(user.Breaches) != 2 {
		t.Fatal("Wrong breach count ", len(user.Breaches), " ", user.Breaches)
	}

	b1Found := false
	b2Found := false

	for _, br := range user.Breaches {
		if br.Title == b.Title {
			b1Found = true
		}
		if br.Title == b2.Title {
			b2Found = true
		}
	}

	if !b1Found {
		t.Fatal("b1 not found")
	}

	if !b2Found {
		t.Fatal("b2 not found")
	}
}

func TestSaveBreach(t *testing.T) {
	c := &config{
		DatabaseDialect: "sqlite3",
		DatabaseArgs:    ":memory:",
	}
	db, err := c.GetDatabase()

	if err != nil {
		t.Fatal(err)
	}

	u := &updater{db: db}

	b1 := &Breach{
		Title:        "title",
		Name:         "name",
		Domain:       "once.tld",
		AddedDate:    time.Now(),
		ModifiedDate: time.Now(),
		Description:  "abc",
	}

	b1clone := &Breach{
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

	br1, err := u.saveBreach(b1)
	if err != nil {
		t.Fatal(err)
	}
	if br1.ID == 0 {
		t.Fatal("No ID for br1")
	}

	br2, err := u.saveBreach(b1clone)
	if err != nil {
		t.Fatal(err)
	}
	if br1.ID != br2.ID {
		t.Fatal("ID missmatch br1 ", br1.ID, " br2 ", br2.ID)
	}

	saved := []Breach{}
	if err := db.Find(&saved).Error; err != nil {
		t.Fatal(err)
	}

	if len(saved) != 1 {
		t.Fatal("Wrong saved breach count ", len(saved), " ", saved)
	}

	br3, err := u.saveBreach(b2)
	if err != nil {
		t.Fatal(err)
	}
	if br3.ID == 0 {
		t.Fatal("No ID for br1")
	}

	saved = []Breach{}
	if err := db.Find(&saved).Error; err != nil {
		t.Fatal(err)
	}

	if len(saved) != 2 {
		t.Fatal("Wrong saved breach count ", len(saved), " ", saved)
	}
}

func TestHandleSource(t *testing.T) {
	c := &config{
		DatabaseDialect: "sqlite3",
		DatabaseArgs:    ":memory:",
	}
	db, err := c.GetDatabase()

	if err != nil {
		t.Fatal(err)
	}

	hibp := &hibpClientMock{}
	ds := newMockDataSource([]string{accountBreachedOnce, accountBreachedTwice, accountBreachedNever})

	u := &updater{
		db:          db,
		newBreached: []string{},
		hibp:        hibp,
	}

	err = u.handleSource(ds)

	if err != nil {
		t.Fatal(err)
	}

	var userOnce User
	if err := db.Where("email = ?", accountBreachedOnce).Preload("Breaches").First(&userOnce).Error; err != nil {
		t.Fatal(err)
	}
	if len(userOnce.Breaches) != 1 {
		t.Fatal("Wrong breach count ", len(userOnce.Breaches), " ", userOnce.Breaches)
	}

	var userTwice User
	if err := db.Where("email = ?", accountBreachedTwice).Preload("Breaches").First(&userTwice).Error; err != nil {
		t.Fatal(err)
	}
	if len(userTwice.Breaches) != 2 {
		t.Fatal("Wrong breach count ", len(userTwice.Breaches), " ", userTwice.Breaches)
	}

	var userNever User
	if err := db.Where("email = ?", accountBreachedNever).Preload("Breaches").First(&userNever).Error; err != gorm.ErrRecordNotFound {
		t.Fatal("Wrong error - user should not exist ", err)
	}
}

func TestNewUpdate(t *testing.T) {
	c := &config{
		DatabaseDialect: "sqlite3",
		DatabaseArgs:    ":memory:",
	}

	hibp := &hibpClientMock{}
	ds := []dataSource{
		newMockDataSource([]string{accountBreachedOnce, accountBreachedNever}),
		newMockDataSource([]string{accountBreachedOnce, accountBreachedTwice}),
	}

	n := &notificationMock{}

	u, err := newUpdater(c, []notification{n}, ds, hibp)

	if err != nil {
		t.Fatal(err)
	}

	u.update()

	var userOnce User
	if err := u.db.Where("email = ?", accountBreachedOnce).Preload("Breaches").First(&userOnce).Error; err != nil {
		t.Fatal(err)
	}
	if len(userOnce.Breaches) != 1 {
		t.Fatal("Wrong breach count ", len(userOnce.Breaches), " ", userOnce.Breaches)
	}

	var userTwice User
	if err := u.db.Where("email = ?", accountBreachedTwice).Preload("Breaches").First(&userTwice).Error; err != nil {
		t.Fatal(err)
	}
	if len(userTwice.Breaches) != 2 {
		t.Fatal("Wrong breach count ", len(userTwice.Breaches), " ", userTwice.Breaches)
	}

	var userNever User
	if err := u.db.Where("email = ?", accountBreachedNever).Preload("Breaches").First(&userNever).Error; err != gorm.ErrRecordNotFound {
		t.Fatal("Wrong error - user should not exist ", err)
	}

	if len(n.reported) != 2 {
		t.Fatal("Wrong reported count ", len(n.reported), " ", n.reported)
	}
}

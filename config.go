package hibpnotify

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const HIBP_API_URL = "https://haveibeenpwned.com/api/v2/breachedaccount/"

type config struct {
	DatabaseDialect string
	DatabaseArgs    string

	CheckIntervalInSeconds int // frequency in which to check for new breaches (min: 3600 s== 1 hour)

	NotifyEmail        string // email address to send notifications about newly detected breaches to
	NotifySMTPHost     string
	NotifySMTPAddr     string
	NotifySMTPUser     string
	NotifySMTPPassword string
	NotifySMTPFrom     string

	CSVPath string

	db *gorm.DB
}

func NewConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var c config
	err = json.Unmarshal(data, &c)
	return &c, err
}

func InitConfig(path string) error {
	c := &config{
		CheckIntervalInSeconds: 3600 * 12, // default: every 12 hours
		DatabaseDialect:        "sqlite3",
		DatabaseArgs:           "hibpnotify.sqlite3",
		CSVPath:                "emails.csv",
		NotifyEmail:            "ops@company.tld",
		NotifySMTPHost:         "mail.company.tld",
		NotifySMTPAddr:         "mail.company.tld:587",
		NotifySMTPUser:         "smptuser",
		NotifySMTPPassword:     "smtppass",
		NotifySMTPFrom:         "youHaveBeenPwnd@company.tld",
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

func (c *config) HIBPBaseURL() string {
	return HIBP_API_URL
}

func (c *config) GetDatabase() (*gorm.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	db, err := gorm.Open(c.DatabaseDialect, c.DatabaseArgs)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Breach{})
	db.AutoMigrate(&User{})

	c.db = db
	return db, nil
}

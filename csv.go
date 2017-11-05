package hibpnotify

import (
	"encoding/csv"
	"log"
	"os"
)

type csvSource struct {
	path string
}

func newCSVSource(c *config) (dataSource, error) {
	s := csvSource{
		path: c.CSVPath,
	}
	return &s, nil
}

func (c *csvSource) GetAccounts() ([]string, error) {
	accounts := []string{}
	csvf, err := os.Open(c.path)

	if err != nil {
		return nil, err
	}

	defer csvf.Close()

	reader := csv.NewReader(csvf)
	fields, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	for i, email := range fields {
		if i == 0 { // skip header
			continue
		}
		accounts = append(accounts, email[0])
	}

	log.Println("Sources: ", accounts)

	return accounts, nil
}

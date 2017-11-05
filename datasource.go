package hibpnotify

type dataSource interface {
	GetAccounts() ([]string, error)
}

func newDataSources(c *config) ([]dataSource, error) {
	ds := []dataSource{}

	if c.CSVPath != "" {
		csv, err := newCSVSource(c)

		if err != nil {
			return nil, err
		}

		ds = append(ds, csv)
	}

	return ds, nil
}

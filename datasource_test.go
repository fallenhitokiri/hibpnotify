package hibpnotify

type dataSourceMock struct {
	accounts []string
}

func newMockDataSource(accounts []string) *dataSourceMock {
	d := &dataSourceMock{
		accounts: accounts,
	}
	return d
}

func (d *dataSourceMock) GetAccounts() ([]string, error) {
	return d.accounts, nil
}

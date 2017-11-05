package hibpnotify

type HIBPNotify struct {
	Config      *config
	DataSources []dataSource
	HIBPClient  hibpClient
	Updater     *updater
}

func New(cfgPath string) (*HIBPNotify, error) {
	hibpn := &HIBPNotify{}

	c, err := NewConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	n, err := newNotification(c)
	if err != nil {
		return nil, err
	}

	ds, err := newDataSources(c)
	if err != nil {
		return nil, err
	}

	hibp, err := newHIBPClient(c)
	if err != nil {
		return nil, err
	}

	u, err := newUpdater(c, n, ds, hibp)
	if err != nil {
		return nil, err
	}

	hibpn.Config = c
	hibpn.DataSources = ds
	hibpn.HIBPClient = hibp
	hibpn.Updater = u

	return hibpn, nil
}

func (h *HIBPNotify) Run() {
	for {
		h.Updater.update()
	}
}

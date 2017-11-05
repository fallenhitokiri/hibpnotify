package hibpnotify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type hibpClient interface {
	requestForEmail(email string) ([]*Breach, error)
}

type hibpApiClient struct {
	baseURL string
	client  *http.Client
}

func newHIBPClient(c *config) (hibpClient, error) {
	hibp := &hibpApiClient{
		baseURL: c.HIBPBaseURL(),
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
	return hibp, nil
}

func (c *hibpApiClient) requestForEmail(email string) ([]*Breach, error) {
	uri := c.baseURL + email
	response, err := c.client.Get(uri)

	if err != nil {
		return nil, err
	}

	// check response codes

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var breaches []*Breach

	err = json.Unmarshal(data, &breaches)

	if err != nil {
		return nil, err
	}

	return breaches, nil
}

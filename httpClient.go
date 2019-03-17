package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Quote struct {
	Subject string `json:"subject"`
	Author  string `json:"author"`
	Text    string `json:"text"`
}

type HttpQuoteFetcher struct {
	hostname string
}

func New(hostname string) *HttpQuoteFetcher {
	return &HttpQuoteFetcher{hostname: hostname}
}

func (fetcher HttpQuoteFetcher) FetchQuote(subject string) (*Quote, error) {
	response, err := http.Get(fetcher.hostname + "/api/" + subject)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	quote := &Quote{}
	if err := json.Unmarshal(respBody, &quote); err != nil {
		return nil, err
	}

	return quote, nil
}

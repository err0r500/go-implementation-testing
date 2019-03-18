package client

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Quote struct {
	Subject string `json:"subject"`
	Author  string `json:"author"`
	Text    string `json:"text"`
}

type HttpQuoteFetcher struct {
	hostname      string
	httpGet       HttpGetter
	readBody      BodyReader
	unMarshalResp RespUnmarshaller
}
type HttpGetter func(string) (*http.Response, error)
type BodyReader func(io.Reader) ([]byte, error)
type RespUnmarshaller func([]byte, interface{}) error

func New(hostname string, mutators ...func(*HttpQuoteFetcher)) *HttpQuoteFetcher {
	fetch := &HttpQuoteFetcher{
		hostname:      hostname,
		httpGet:       http.Get,
		readBody:      ioutil.ReadAll,
		unMarshalResp: json.Unmarshal,
	}

	for _, mutator := range mutators {
		mutator(fetch)
	}
	return fetch
}

func SetHttpGet(customHttpGet HttpGetter) func(*HttpQuoteFetcher) {
	return func(fetcher *HttpQuoteFetcher) {
		fetcher.httpGet = customHttpGet
	}
}
func SetReadBody(bodyReader BodyReader) func(*HttpQuoteFetcher) {
	return func(fetcher *HttpQuoteFetcher) {
		fetcher.readBody = bodyReader
	}
}

func SetRespUnmarshaller(respUnmarshaller RespUnmarshaller) func(*HttpQuoteFetcher) {
	return func(fetcher *HttpQuoteFetcher) {
		fetcher.unMarshalResp = respUnmarshaller
	}
}

func (fetcher HttpQuoteFetcher) FetchQuote(subject string) (*Quote, error) {
	response, err := fetcher.httpGet(fetcher.hostname + "/api/" + subject)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("httpGet returned a nil pointer")
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	respBody, err := fetcher.readBody(response.Body)
	if err != nil {
		return nil, err
	}

	quote := &Quote{}
	if err := fetcher.unMarshalResp(respBody, &quote); err != nil {
		return nil, err
	}

	return quote, nil
}

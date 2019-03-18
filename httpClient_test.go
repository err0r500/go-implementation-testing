package client_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	client "github.com/err0r500/go-implementation-testing"

	"github.com/stretchr/testify/assert"
)

const subject = "art"

var validQuote = &client.Quote{Author: "author", Text: "text", Subject: subject}

func TestHttpQuoteFetcher_FetchQuote(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/"+subject, r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				validResponse(w, validQuote)
			}),
		)
		defer fakeServer.Close()

		quote, err := client.New(fakeServer.URL).FetchQuote(subject)
		assert.NoError(t, err)
		assert.Equal(t, validQuote, quote)
	})

	t.Run("response status error", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", http.StatusBadRequest)
			}),
		)
		defer fakeServer.Close()

		_, err := client.New(fakeServer.URL).FetchQuote(subject)
		assert.Error(t, err)
	})

	t.Run("more tricky test cases", func(t *testing.T) {
		t.Run("error on httpGet", func(t *testing.T) {
			failureTest(t,
				client.SetHttpGet(func(string) (*http.Response, error) {
					return nil, errors.New("")
				}),
			)
		})

		t.Run("nil on httpGet", func(t *testing.T) {
			failureTest(t,
				client.SetHttpGet(func(string) (*http.Response, error) {
					return nil, nil
				}),
			)
		})

		t.Run("error on readBody", func(t *testing.T) {
			failureTest(t,
				client.SetReadBody(func(io.Reader) ([]byte, error) {
					return nil, errors.New("")
				}),
			)
		})

		t.Run("error on unMarshalResp", func(t *testing.T) {
			failureTest(t,
				client.SetRespUnmarshaller(func([]byte, interface{}) error {
					return errors.New("")
				}),
			)
		})
	})
}

func failureTest(t *testing.T, mutator func(*client.HttpQuoteFetcher)) {
	fakeServer := getDummyServer()
	defer fakeServer.Close()
	_, err := client.New(
		fakeServer.URL,
		mutator,
	).FetchQuote(subject)

	assert.Error(t, err)
}

func validResponse(w http.ResponseWriter, validQuote *client.Quote) {
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(validQuote)
	_, _ = w.Write(resp)
}

func getDummyServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validResponse(w, validQuote)
		}),
	)
}

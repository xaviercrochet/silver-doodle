package clients

import (
	"fmt"
	"net/http"
)

// RestClient ...
var RestClient restClientInterface

type restClientInterface interface {
	Get(url string, header http.Header) (*http.Response, error)
}

type restClient struct{}

func init() {
	RestClient = &restClient{}
}

func (rc *restClient) Get(url string, header http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("cannot send http request: %v", err)
		return nil, err
	}

	request.Header = header
	client := http.Client{}

	return client.Do(request)
}

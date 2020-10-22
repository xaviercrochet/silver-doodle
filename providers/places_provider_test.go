package providers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"localsearch-api/clients"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type restClientMock struct{}

var (
	funcGet func(url string, header http.Header) (*http.Response, error)
)

func (rc *restClientMock) Get(url string, header http.Header) (*http.Response, error) {
	return funcGet(url, header)
}

func TestMain(m *testing.M) {
	clients.RestClient = &restClientMock{}
}

func TestGetPlaceRestClientError(t *testing.T) {
	t.Log("TestGetPlaceRestClientError")
	restError := errors.New("rest client failed")

	funcGet = func(url string, header http.Header) (*http.Response, error) {
		return nil, restError
	}

	place, err := PlacesProvider.GetPlace("123")

	require.NotNil(t, err)
	require.Nil(t, place)
	require.Equal(t, err.StatusCode, http.StatusInternalServerError)
	require.Equal(t, err.Message, fmt.Sprintf("could not retrieve place from api: %v", restError))
}

func TestGetPlaceSuccess(t *testing.T) {
	funcGet = func(url string, header http.Header) (*http.Response, error) {
		json := `{}`
		r := ioutil.NopCloser(strings.NewReader(json))
		return &http.Response{
			Body:       r,
			StatusCode: http.StatusOK,
		}, nil
	}

	place, err := PlacesProvider.GetPlace("123")
	require.Nil(t, err)
	require.NotNil(t, place)
}

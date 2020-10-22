package providers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"localsearch-api/clients"
	"net/http"
	"os"
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

func TestGetPlaceRestClientError(t *testing.T) {

	clients.RestClient = &restClientMock{}
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
	clients.RestClient = &restClientMock{}
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

func TestGetPlaceResponseBodyError(t *testing.T) {
	clients.RestClient = &restClientMock{}
	funcGet = func(url string, header http.Header) (*http.Response, error) {
		invalidCloser, _ := os.Open("-abc")
		return &http.Response{
			Body:       invalidCloser,
			StatusCode: http.StatusOK,
		}, nil
	}

	place, err := PlacesProvider.GetPlace("123")

	require.Nil(t, place)
	require.NotNil(t, err)
	require.Equal(t, err.StatusCode, http.StatusInternalServerError)

}

func TestGetPlaceInvalidResponseStatusCode(t *testing.T) {
	clients.RestClient = &restClientMock{}
	funcGet = func(url string, header http.Header) (*http.Response, error) {
		json := `{}`
		r := ioutil.NopCloser(strings.NewReader(json))
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	place, err := PlacesProvider.GetPlace("123")

	require.Nil(t, place)
	require.NotNil(t, err)
	require.Equal(t, 400, err.StatusCode)

}

func TestGetPlaceInvalidSuccessfulResponseBody(t *testing.T) {
	clients.RestClient = &restClientMock{}
	funcGet = func(url string, header http.Header) (*http.Response, error) {
		json := `{`
		r := ioutil.NopCloser(strings.NewReader(json))

		return &http.Response{
			Body:       r,
			StatusCode: http.StatusOK,
		}, nil
	}

	place, err := PlacesProvider.GetPlace("123")

	require.Nil(t, place)
	require.NotNil(t, err)
	require.Equal(t, http.StatusInternalServerError, err.StatusCode)

}

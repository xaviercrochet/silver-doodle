package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localsearch-api/clients"
	"localsearch-api/config"
	"localsearch-api/domain/localsearch"
	"net/http"
)

// PlacesProvider ...
var PlacesProvider placesProviderInterface

type placesProviderInterface interface {
	GetPlace(placeID string) (*localsearch.Place, *localsearch.ErrorResponse)
}

type placesProvider struct{}

func init() {
	PlacesProvider = &placesProvider{}
}

func (p *placesProvider) GetPlace(placeID string) (*localsearch.Place, *localsearch.ErrorResponse) {
	path := fmt.Sprintf("%s/%s", config.GetExternalAPIURl(), placeID)
	// no specific header for now
	header := http.Header{}

	resp, err := clients.RestClient.Get(path, header)
	if err != nil {
		return nil, &localsearch.ErrorResponse{
			Message:    fmt.Sprintf("could not retrieve place from api: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &localsearch.ErrorResponse{
			Message:    fmt.Sprintf("cannot read response body: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		errorResponse := &localsearch.ErrorResponse{}
		// error message are provided in XML :s
		errorResponse.Message = "error"
		errorResponse.StatusCode = resp.StatusCode
		return nil, errorResponse
	}

	localSearchPlace := &localsearch.Place{}
	if err := json.Unmarshal(bytes, localSearchPlace); err != nil {
		return nil, &localsearch.ErrorResponse{
			Message:    fmt.Sprintf("could not deserialize json body: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return localSearchPlace, nil

}

package services

import (
	"localsearch-api/domain/localsearch"
	"localsearch-api/providers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	funcGetPlace func(placeID string) (*localsearch.Place, *localsearch.ErrorResponse)
)

type placesProviderMock struct {
}

func (p *placesProviderMock) GetPlace(placeID string) (*localsearch.Place, *localsearch.ErrorResponse) {
	return funcGetPlace(placeID)
}

func TestGetPlaceApplicationError(t *testing.T) {
	providers.PlacesProvider = &placesProviderMock{}
	funcGetPlace = func(placeID string) (*localsearch.Place, *localsearch.ErrorResponse) {
		return nil, &localsearch.ErrorResponse{
			Message:    "error",
			StatusCode: http.StatusInternalServerError,
		}
	}

	place, err := PlacesService.GetPlace("123")
	require.Nil(t, place)
	require.NotNil(t, err)
	require.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

func TestGetPlaceSuccess(t *testing.T) {
	providers.PlacesProvider = &placesProviderMock{}

	funcGetPlace = func(placeID string) (*localsearch.Place, *localsearch.ErrorResponse) {
		return &localsearch.Place{
			Addresses: []*localsearch.Address{
				&localsearch.Address{
					Where: &localsearch.Where{
						Street:      "street",
						City:        "city",
						HouseNumber: "house number",
						Zipcode:     float64(0),
					},
				},
			},
			OpeningHours: &localsearch.OpeningHours{
				Days: &localsearch.Days{
					Monday: []*localsearch.Schedule{
						&localsearch.Schedule{
							Start: "10:00",
							End:   "12:00",
							Type:  "OPEN",
						},
						&localsearch.Schedule{
							Start: "14:00",
							End:   "22:00",
							Type:  "OPEN",
						},
					},
				},
			},
		}, nil
	}

	place, err := PlacesService.GetPlace("123")

	require.Nil(t, err)
	require.NotNil(t, place)
	require.Equal(t, "123", place.ID)
	require.NotEmpty(t, place.Location)
}

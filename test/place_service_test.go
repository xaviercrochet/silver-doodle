package service_test

import (
	"localsearch-api/service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPlace(t *testing.T) {
	t.Parallel()
	placeID1 := "GXvPAor1ifNfpF0U5PTG0w"
	placeID2 := "ohGSnJtMIC5nPfYRi_HTAg"
	place, err := service.GetPlace(placeID1)
	require.NoError(t, err)
	require.NotNil(t, place)

	place, err = service.GetPlace(placeID2)
	require.NoError(t, err)
	require.NotNil(t, place)
}

func TestGetPlaceUnknownID(t *testing.T) {
	t.Parallel()
	placeID := "xavier"
	place, err := service.GetPlace(placeID)
	require.NoError(t, err)
	require.Nil(t, place)

}

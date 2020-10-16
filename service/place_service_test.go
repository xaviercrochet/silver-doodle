package service_test

import (
	"localsearch-api/service"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPlace(t *testing.T) {
	t.Parallel()
	placeID1 := "GXvPAor1ifNfpF0U5PTG0w"
	placeID2 := "ohGSnJtMIC5nPfYRi_HTAg"
	place, err := service.GetPlace(placeID1)
	require.Nil(t, err)
	require.NotNil(t, place)
	days := []string{}

	for _, schedule := range place.Schedules {
		require.NotEqual(t, len(schedule.Days), 0)
		require.NotNil(t, schedule.HoursRanges)
		for _, day := range schedule.Days {
			require.NotNil(t, day)
			require.NotEmpty(t, day)
			days = append(days, day)

		}
	}
	require.Equal(t, len(days), 7)

	place, err = service.GetPlace(placeID2)
	require.Nil(t, err)
	require.NotNil(t, place)
	days = []string{}

	for _, schedule := range place.Schedules {
		require.NotEqual(t, len(schedule.Days), 0)
		require.NotNil(t, schedule.HoursRanges)
		for _, day := range schedule.Days {
			require.NotNil(t, day)
			require.NotEmpty(t, day)
			days = append(days, day)

		}
	}
	require.Equal(t, len(days), 7)
}

func TestGetPlaceUnknownID(t *testing.T) {
	t.Parallel()
	placeID := "xavier"
	place, err := service.GetPlace(placeID)
	require.NotNil(t, err)
	require.Equal(t, err.StatusCode, http.StatusNotFound)
	require.Nil(t, place)

}

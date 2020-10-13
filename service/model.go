package service

import (
	"fmt"
	"localsearch-api/json"
	"log"
)

// Place is the place retrieved from the API
type Place struct {
	ID           string
	Name         string
	Location     string
	OpeningHours []*OpeningHour
}

// OpeningHour is a place's opening hour retrieved from the api
type OpeningHour struct {
	Days  []string
	Hours []string
}

// NewOpeningHour returns a new NewOpeningHour with isOpen initialized to false
func NewOpeningHour(Day string) *OpeningHour {
	return &OpeningHour{
		Days:  []string{Day},
		Hours: []string{},
	}
}

func toPlace(localSearchPlace *json.LocalSearchPlace) *Place {
	place := &Place{}
	place.Name = localSearchPlace.DisplayedWhat
	if len(localSearchPlace.Addresses) > 0 {
		place.Location = fmt.Sprintf(
			"%s %s, %d %s",
			localSearchPlace.Addresses[0].Where.Street,
			localSearchPlace.Addresses[0].Where.HouseNumber,
			uint32(localSearchPlace.Addresses[0].Where.Zipcode),
			localSearchPlace.Addresses[0].Where.City,
		)
	}
	log.Print(localSearchPlace.OpeningHours.Days.Tuesday)
	for key, day := range localSearchPlace.OpeningHours.Days.Tuesday {
		log.Printf("Day: %d %v", key, day)
	}

	return place
}

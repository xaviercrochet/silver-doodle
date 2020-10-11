package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

const apiURL = "https://storage.googleapis.com/coding-session-rest-api"

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

//NewWeek return an array of OpeningHour for all days of the week
func NewWeek() map[string]*OpeningHour {
	result := make(map[string]*OpeningHour)

	result["monday"] = NewOpeningHour("monday")
	result["tuesday"] = NewOpeningHour("tuesday")
	result["wednesday"] = NewOpeningHour("wednesday")
	result["thursday"] = NewOpeningHour("thursday")
	result["friday"] = NewOpeningHour("friday")
	result["saturday"] = NewOpeningHour("saturday")
	result["sunday"] = NewOpeningHour("sunday")

	return result
}

func (place *Place) computeSchedule(openingHours map[string]*OpeningHour) {
	rawSchedule := []*OpeningHour{
		openingHours["monday"],
		openingHours["tuesday"],
		openingHours["wednesday"],
		openingHours["thursday"],
		openingHours["friday"],
		openingHours["saturday"],
		openingHours["sunday"],
	}

	for _, day := range rawSchedule {
		if len(place.OpeningHours) == 0 {
			place.OpeningHours = []*OpeningHour{day}

		} else {
			lastIndex := len(place.OpeningHours) - 1

			if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, day.Hours) {
				place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, day.Days[0])
			} else {
				place.OpeningHours = append(place.OpeningHours, day)
			}
		}
	}
}

// GetPlace retrieve a Place by the provided ID
func GetPlace(placeID string) (*Place, error) {
	// todo checkplaceID is a valid one
	path := fmt.Sprintf("%s/%s", apiURL, placeID)
	log.Printf("GetPath: %v", path)

	resp, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve place from api: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Printf("expected 200 http status code but got in stead %v", resp.StatusCode)
		return nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse json body: %v", err)
	}

	place, err := ParsePlaceJSON(body)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize json body: %v", err)
	}
	place.ID = placeID
	return place, nil
}

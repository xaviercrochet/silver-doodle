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
	monday := openingHours["monday"]
	tuesday := openingHours["tuesday"]
	wednesday := openingHours["wednesday"]
	thursday := openingHours["thursday"]
	friday := openingHours["friday"]
	saturday := openingHours["saturday"]
	sunday := openingHours["sunday"]

	//monday
	place.OpeningHours = []*OpeningHour{monday}

	//tuesday
	lastIndex := len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, tuesday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, tuesday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, tuesday)
	}

	//wednesday
	lastIndex = len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, wednesday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, wednesday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, wednesday)
	}

	//thursday
	lastIndex = len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, thursday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, thursday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, thursday)
	}

	//friday
	lastIndex = len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, friday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, friday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, friday)
	}

	//saturday
	lastIndex = len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, saturday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, saturday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, saturday)
	}

	//sunday
	lastIndex = len(place.OpeningHours) - 1
	if reflect.DeepEqual(place.OpeningHours[lastIndex].Hours, sunday.Hours) {
		place.OpeningHours[lastIndex].Days = append(place.OpeningHours[lastIndex].Days, sunday.Days[0])
	} else {
		place.OpeningHours = append(place.OpeningHours, sunday)
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

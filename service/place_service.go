package service

import (
	"fmt"
	"io/ioutil"
	"localsearch-api/json"
	"log"
	"net/http"
	"reflect"
)

const apiURL = "https://storage.googleapis.com/coding-session-rest-api"

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
		return nil, logError(fmt.Errorf("could not retrieve place from api: %v", err))

	}

	if resp.StatusCode != 200 {
		return nil, logError(fmt.Errorf("expected http 200 but got %v", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, logError(fmt.Errorf("could not parse json body: %v", err))
	}

	//place, err := ParsePlaceJSON(body)
	rawPlace, err := json.Parse(body)
	if err != nil {
		return nil, logError(fmt.Errorf("could not deserialize json body: %v", err))

	}
	log.Printf("Raw Response: %v", rawPlace)
	place := toPlace(rawPlace)

	place.ID = placeID
	return place, nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

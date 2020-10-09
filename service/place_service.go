package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

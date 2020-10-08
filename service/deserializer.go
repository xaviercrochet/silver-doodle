package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ParsePlaceJSON parse Place data out of body and returns a new Place
func ParsePlaceJSON(body []byte) (*Place, error) {
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	//fmt.Print(result)

	place := &Place{
		OpeningHours: make(map[string]*OpeningHour),
	}

	for key, value := range result {
		switch key {
		case "displayed_what":
			place.Name = value.(string)
		case "addresses":
			locations := []string{}
			for _, value := range value.([]interface{}) {
				address := value.(map[string]interface{})
				where := address["where"].(map[string]interface{})
				street := where["street"].(string)
				city := where["city"].(string)
				houseNumber := where["house_number"].(string)
				zipcode := where["zipcode"].(float64)
				location := fmt.Sprintf("%s %s, %d %s", street, houseNumber, uint32(zipcode), city)
				locations = append(locations, location)
			}
			if len(locations) > 0 {
				place.Location = locations[0]
			} else {
				place.Location = "undefined"
			}
		case "opening_hours":
			openingHourJSON := value.(map[string]interface{})
			daysJSON := openingHourJSON["days"].(map[string]interface{})
			for key, value := range daysJSON {
				openingHour := &OpeningHour{}
				for _, value := range value.([]interface{}) {
					hoursJSON := value.(map[string]interface{})
					openingType := hoursJSON["type"].(string)
					start := hoursJSON["start"].(string)
					end := hoursJSON["end"].(string)
					if openingType == "OPEN" {
						openingValue := fmt.Sprintf("%s - %s", start, end)
						openingHour.Values = append(openingHour.Values, openingValue)
					} else {
						openingHour.Values = []string{openingType}
						break
					}
				}
				// check for duplicates
				found := false
				for day, existingOpeningHour := range place.OpeningHours {
					if reflect.DeepEqual(existingOpeningHour.Values, openingHour.Values) {
						splits := strings.Split(day, " - ")

						if len(splits) > 0 {
							days := fmt.Sprintf("%s - %s", splits[0], key)
							place.OpeningHours[days] = openingHour

						} else {
							days := fmt.Sprintf("%s - %s", day, key)
							place.OpeningHours[days] = openingHour

						}
						delete(place.OpeningHours, day)
						found = true
						break
					}
				}
				if !found {
					place.OpeningHours[key] = openingHour
				}
			}

		default:
			//fmt.Println(key)
		}
	}

	return place, nil
}

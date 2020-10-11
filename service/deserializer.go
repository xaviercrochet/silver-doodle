package service

import (
	"encoding/json"
	"fmt"
)

// ParsePlaceJSON parse Place data out of body and returns a new Place
func ParsePlaceJSON(body []byte) (*Place, error) {
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	place := &Place{}

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
			openingHours := NewWeek()

			for day, value := range daysJSON {
				openingHour := openingHours[day]
				if openingHour != nil {
					for _, value := range value.([]interface{}) {
						hoursJSON := value.(map[string]interface{})
						openingType := hoursJSON["type"].(string)
						start := hoursJSON["start"].(string)
						end := hoursJSON["end"].(string)
						if openingType == "OPEN" {
							openingValue := fmt.Sprintf("%s - %s", start, end)
							openingHour.Hours = append(openingHour.Hours, openingValue)
						}
					}
				}
			}
			place.computeSchedule(openingHours)

		default:
			//fmt.Println(key)
		}
	}

	return place, nil
}

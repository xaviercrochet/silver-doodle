package service

import (
	"fmt"
	"io/ioutil"
	"localsearch-api/json"
	"localsearch-api/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiURL = "https://storage.googleapis.com/coding-session-rest-api"

// GetPlace retrieve a Place by the provided ID
func GetPlace(placeID string) (*Place, *utils.ApplicationError) {
	// todo checkplaceID is a valid one
	path := fmt.Sprintf("%s/%s", apiURL, placeID)
	log.Printf("GetPath: %v", path)

	resp, err := http.Get(path)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    fmt.Sprintf("could not retrieve place from api: %v", err),
			StatusCode: http.StatusServiceUnavailable,
			Code:       "localsearch api not available",
		}
		return nil, apiErr

	}

	if resp.StatusCode != 200 {
		apiErr := &utils.ApplicationError{
			Message:    fmt.Sprintf("could not retrieve place from api: %v", resp.StatusCode),
			StatusCode: http.StatusNotFound,
			Code:       "resource not found",
		}
		return nil, apiErr
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    fmt.Sprintf("could not parse json body: %v", err),
			StatusCode: http.StatusServiceUnavailable,
			Code:       "local search api not available",
		}
		return nil, apiErr
	}

	rawPlace, err := json.Parse(body)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    fmt.Sprintf("could not deserialize json body: %v", err),
			StatusCode: http.StatusInternalServerError,
			Code:       "internal server error",
		}
		return nil, apiErr

	}
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
	monday := parseDaySchedules(time.Weekday(time.Monday), localSearchPlace.OpeningHours.Days.Monday)
	tuesday := parseDaySchedules(time.Weekday(time.Tuesday), localSearchPlace.OpeningHours.Days.Tuesday)
	wednesday := parseDaySchedules(time.Weekday(time.Wednesday), localSearchPlace.OpeningHours.Days.Wednesday)
	thursday := parseDaySchedules(time.Weekday(time.Thursday), localSearchPlace.OpeningHours.Days.Thursday)
	friday := parseDaySchedules(time.Weekday(time.Friday), localSearchPlace.OpeningHours.Days.Friday)
	saturday := parseDaySchedules(time.Weekday(time.Saturday), localSearchPlace.OpeningHours.Days.Saturday)
	sunday := parseDaySchedules(time.Weekday(time.Sunday), localSearchPlace.OpeningHours.Days.Sunday)

	place.Schedules = []*Schedule{
		monday,
		tuesday,
		wednesday,
		thursday,
		friday,
		saturday,
		sunday,
	}

	if !place.isOpenToday() {
		place.OpenNext = findOpenNext(place.Schedules)
	}

	return place
}

func parseDaySchedules(weekday time.Weekday, localSearchOpeningHours []*json.LocalSearchSchedule) *Schedule {
	schedule := NewSchedule(weekday.String())
	schedule.Days = []string{weekday.String()}
	for _, localSearchOpeningHour := range localSearchOpeningHours {
		hoursRange := NewHoursRange(localSearchOpeningHour.Start, localSearchOpeningHour.End)
		schedule.HoursRanges = append(schedule.HoursRanges, hoursRange)
		hoursRange.IsOpenNow = localSearchOpeningHour.IsOpen(weekday)

	}
	return schedule
}

func computeDate(current time.Time, hourAndMins string) time.Time {

	splits := strings.Split(hourAndMins, ":")
	hour, err := strconv.Atoi(splits[0])
	if err != nil {
		log.Printf("could not parse hour part of schedule.Start: %v", err)
		return current
	}
	min, err := strconv.Atoi(splits[1])
	if err != nil {
		log.Printf("could not parse minuts part of schedule.Start: %v", err)
		return current
	}

	return time.Date(current.Year(), current.Month(), current.Day(), hour, min, 0, 0, current.Location())
}

func findOpenNext(schedules []*Schedule) *Schedule {
	//1 find index of today
	now := time.Now()
	index := 0
	for i, schedule := range schedules {
		if schedule.Days[0] == now.Weekday().String() {
			index = i
			break
		}
	}
	// 2 look forward

	for i := index; i < len(schedules); i++ {
		schedule := schedules[i]
		for j, hoursRange := range schedule.HoursRanges {
			start := computeDate(now, hoursRange.Start)
			end := computeDate(now, hoursRange.End)
			if now.After(start) && now.Before(end) {
				return schedule
			} else if j < len(schedule.HoursRanges)-1 {
				return schedule
			}
		}
		now.AddDate(0, 0, 1)
	}
	// 3 look backward
	now = time.Now()
	for i := index; i >= 0; i-- {
		schedule := schedules[i]
		for j, hoursRange := range schedule.HoursRanges {
			start := computeDate(now, hoursRange.Start)
			end := computeDate(now, hoursRange.End)
			if now.After(start) && now.Before(end) {
				return schedule
			} else if j < len(schedule.HoursRanges)-1 {
				return schedule
			}
		}
		now.AddDate(0, 0, -1)

	}
	return nil
}

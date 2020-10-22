package services

import (
	"fmt"
	"localsearch-api/domain/localsearch"
	"localsearch-api/domain/places"
	"localsearch-api/providers"
	"localsearch-api/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

// PlacesService ...
var PlacesService placesServiceInterface

type placesServiceInterface interface {
	GetPlace(placeID string) (*places.Place, *utils.ApplicationError)
}

type placesService struct{}

func init() {
	PlacesService = &placesService{}
}

func (s *placesService) GetPlace(placeID string) (*places.Place, *utils.ApplicationError) {
	resp, err := providers.PlacesProvider.GetPlace(placeID)
	if err != nil {
		return nil, &utils.ApplicationError{
			Message:    err.Message,
			StatusCode: err.StatusCode,
		}
	}
	place := buildPlace(resp)
	place.ID = placeID
	return place, nil
}

func buildPlace(localSearchPlace *localsearch.LocalSearchPlace) *places.Place {
	place := &places.Place{
		Name: localSearchPlace.DisplayedWhat,
	}

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

	place.Schedules = []*places.Schedule{
		monday,
		tuesday,
		wednesday,
		thursday,
		friday,
		saturday,
		sunday,
	}

	if !place.IsOpenToday() {
		place.OpenNext = findOpenNext(place.Schedules)
	}

	return place
}

func parseDaySchedules(weekday time.Weekday, localSearchOpeningHours []*localsearch.LocalSearchSchedule) *places.Schedule {
	schedule := places.NewSchedule(weekday.String())
	schedule.Days = []string{weekday.String()}
	for _, localSearchOpeningHour := range localSearchOpeningHours {
		hoursRange := places.NewHoursRange(localSearchOpeningHour.Start, localSearchOpeningHour.End)
		schedule.HoursRanges = append(schedule.HoursRanges, hoursRange)
		hoursRange.IsOpenNow = localSearchOpeningHour.IsOpen(weekday)

	}
	return schedule
}

func findOpenNext(schedules []*places.Schedule) *places.Schedule {
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

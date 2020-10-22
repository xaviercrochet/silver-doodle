package localsearch

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// Place ...
type Place struct {
	DisplayedWhat string `json:"displayed_what"`
	Addresses     []*Address
	OpeningHours  *OpeningHours `json:"opening_hours"`
}

// OpeningHours ...
type OpeningHours struct {
	Days *Days
}

// Days ...
type Days struct {
	Monday    []*Schedule
	Tuesday   []*Schedule
	Wednesday []*Schedule
	Thursday  []*Schedule
	Friday    []*Schedule
	Saturday  []*Schedule
	Sunday    []*Schedule
}

// Schedule ...
type Schedule struct {
	Start string
	End   string
	Type  string
}

// Address ...
type Address struct {
	Where *Where
}

// Where ...
type Where struct {
	Street      string
	City        string
	HouseNumber string `json:"house_number"`
	Zipcode     float64
}

// IsOpen return true if now is between schedule.start, schedule.end hours and weekday for the given date
func (schedule *Schedule) IsOpen(weekday time.Weekday) bool {
	now := time.Now()
	if now.Weekday() != weekday {
		return false
	}
	startSplits := strings.Split(schedule.Start, ":")
	endSplits := strings.Split(schedule.End, ":")
	startHour, err := strconv.Atoi(startSplits[0])
	if err != nil {
		log.Printf("could not parse hour part of schedule.Start: %v", err)
		return false
	}
	startMin, err := strconv.Atoi(startSplits[1])
	if err != nil {
		log.Printf("could not parse minuts part of schedule.Start: %v", err)
		return false
	}
	endHour, err := strconv.Atoi(endSplits[0])
	if err != nil {
		log.Printf("could not parse hour part of schedule.End: %v", err)
		return false
	}
	endMin, err := strconv.Atoi(endSplits[1])
	if err != nil {
		log.Printf("could not parse minuts part of schedule.End: %v", err)
		return false
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, now.Location())

	return now.After(start) && now.Before(end)
}

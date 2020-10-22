package places

// Place is the place retrieved from the API
type Place struct {
	ID        string
	Name      string
	Location  string
	Schedules []*Schedule
	OpenNext  *Schedule
}

// Schedule is a place's opening hour retrieved from the api
type Schedule struct {
	Days        []string
	HoursRanges []*HoursRange
}

// HoursRange ...
type HoursRange struct {
	Start     string
	End       string
	IsOpenNow bool
}

// NewHoursRange returns a new HoursRange
func NewHoursRange(start string, end string) *HoursRange {
	return &HoursRange{
		Start:     start,
		End:       end,
		IsOpenNow: false,
	}
}

// NewSchedule returns a new NewSchedule
func NewSchedule(Day string) *Schedule {
	return &Schedule{
		Days:        []string{Day},
		HoursRanges: []*HoursRange{},
	}
}

func (s *Schedule) compareHoursRanges(schedule *Schedule) bool {
	if len(s.HoursRanges) != len(schedule.HoursRanges) {
		return false
	}

	for i := range s.HoursRanges {
		sameStart := s.HoursRanges[i].Start == schedule.HoursRanges[i].Start
		sameEnd := s.HoursRanges[i].End == schedule.HoursRanges[i].End
		if !(sameStart && sameEnd) {
			return false
		}
	}

	return true
}

// IsOpenToday ...
func (p *Place) IsOpenToday() bool {
	for _, schedule := range p.Schedules {
		for _, hoursRange := range schedule.HoursRanges {
			if hoursRange.IsOpenNow {
				return true
			}

		}
	}
	return false
}

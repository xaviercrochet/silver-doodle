package json

// LocalSearchPlace ...
type LocalSearchPlace struct {
	DisplayedWhat string `json:"displayed_what"`
	Addresses     []*LocalSearchAddress
	OpeningHours  *LocalSearchOpeningHours `json:"opening_hours"`
}

// LocalSearchOpeningHours ...
type LocalSearchOpeningHours struct {
	Days LocalSearchDays
}

// LocalSearchDays ...
type LocalSearchDays struct {
	Monday    []*LocalSearchSchedule
	Tuesday   []*LocalSearchSchedule
	Wednesday []*LocalSearchSchedule
	Thursday  []*LocalSearchSchedule
	Friday    []*LocalSearchSchedule
	Saturday  []*LocalSearchSchedule
	Sunday    []*LocalSearchSchedule
}

// LocalSearchSchedule ...
type LocalSearchSchedule struct {
	Start string
	End   string
	Type  string
}

// LocalSearchAddress ...
type LocalSearchAddress struct {
	Where *LocalSearchWhere
}

// LocalSearchWhere ...
type LocalSearchWhere struct {
	Street      string
	City        string
	HouseNumber string `json:"house_number"`
	Zipcode     float64
}

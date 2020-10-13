package json

import "encoding/json"

// Parse deserialize localsearch response into a LocalSearchPlace struct
func Parse(bodyJSON []byte) (*LocalSearchPlace, error) {
	place := &LocalSearchPlace{}
	err := json.Unmarshal([]byte(bodyJSON), &place)
	return place, err

}

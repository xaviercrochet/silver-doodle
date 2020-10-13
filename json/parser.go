package json

import "encoding/json"

// Parse ...
func Parse(bodyJSON []byte) (*LocalSearchPlace, error) {
	place := &LocalSearchPlace{}
	err := json.Unmarshal([]byte(bodyJSON), &place)
	return place, err

}

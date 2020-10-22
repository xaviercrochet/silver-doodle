package config

import (
	"os"
)

const (
	apiEndPoint = "API_END_POINT"
)

var (
	externalAPIUrl = os.Getenv(apiEndPoint)
)

// GetExternalAPIURl ...
func GetExternalAPIURl() string {
	if len(externalAPIUrl) == 0 {
		return "https://storage.googleapis.com/coding-session-rest-api"
	}
	return externalAPIUrl
}

package utils

// ApplicationError represent the error returned by the api
type ApplicationError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
}

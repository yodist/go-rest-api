package models

// Response is struct to handle api response
type Response struct {
	StatusCode int         `json:"status_code"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

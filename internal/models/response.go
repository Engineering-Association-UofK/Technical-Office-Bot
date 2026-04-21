package models

import "time"

type DefaultResponse struct {
	Status    string    `json:"status" doc:"The status of the request (success/error)"`
	Message   string    `json:"message,omitempty" doc:"A descriptive message"`
	Timestamp time.Time `json:"timestamp" doc:"The server time of the response"`
}

type RespWithBody struct {
	Body interface{} `json:"body"`
}

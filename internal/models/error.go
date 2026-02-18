package models

import (
	"encoding/json"
	"time"
)

type Error struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timeStamp"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var aux struct {
		Status    int    `json:"status"`
		Message   string `json:"message"`
		TimeStamp int64  `json:"timeStamp"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Status = aux.Status
	e.Message = aux.Message
	e.TimeStamp = time.Unix(aux.TimeStamp, 0)
	return nil
}

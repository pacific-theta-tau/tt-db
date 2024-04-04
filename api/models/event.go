package models

import("time")

// TODO: decide which fields are nullable.
type Event struct {
	EventID			int			`json:"eventID"`
	Category		int			`json:"category"` 
	EventName		string 		`json:"eventName"`
	EventLocation	string 		`json:"eventLocation"`
	EventDate		time.Time	`json:"eventDate"`
}

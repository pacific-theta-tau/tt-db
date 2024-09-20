package models

import("time")

// TODO: decide which fields are nullable.
type Event struct {
	EventID			int			`json:"eventID"`  //primary
	EventName		string 		`json:"eventName"`
	CategoryName	string      `json:"categoryName"` 
	EventLocation	string 		`json:"eventLocation"`
	EventDate		time.Time	`json:"eventDate"`
}

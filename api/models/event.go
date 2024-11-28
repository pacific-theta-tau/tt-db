package models

import("time")

//  @Description Event information
type Event struct {
	EventID			int			`json:"eventID"`  //primary
	EventName		string 		`json:"eventName"`
	CategoryName	string      `json:"categoryName"` 
	EventLocation	string 		`json:"eventLocation"`
	EventDate		time.Time	`json:"eventDate"`
}

//  @Description Event Attendance information of a Brother
type EventAttendance struct {
    BrotherID           int `json:"brotherID"`
    FirstName           string `json:"firstName"`
    LastName            string `json:"lastName"`
    RollCall            int `json:"rollCall"`
    AttendanceStatus    string `json:"attendanceStatus"`
    EventID             int `json:"eventID"`
}

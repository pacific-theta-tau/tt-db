package models

import("time")


var AttendanceStatus = map[string]bool{
    "Present": true,
    "Absent": true,
    "Excused": true,
}

type Attendance struct {
    BrotherID         int    `json:"brotherID"`
    EventID           int    `json:"eventID"`
    AttendanceStatus  string `json:"attendanceStatus"`
    //
    RollCall          int    `json:"rollCall"`
    FirstName         string `json:"firstName"`
    LastName          string `json:"lastName"`
    //
    EventName         string `json:"eventName"`
    EventLocation     string `json:"eventLocation"`
    EventDate         time.Time `json:"eventDate"`        
    EventCategory     string    `json:"eventCategory"`
}

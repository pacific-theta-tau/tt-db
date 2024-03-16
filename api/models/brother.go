package models

// TODO: decide which fields are nullable.
type Brother struct {
	PacificId   string `json:"pacificId"` // Required
	FirstName   string `json:"firstName"` // Required
	LastName    string `json:"lastName"`  // Required
	Status      string `json:"status"`    // Required
	Class       string `json:"className"`
	RollCall    string `json:"rollCall"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	BadStanding int    `json:"badStanding"`
}

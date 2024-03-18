package models

// TODO: decide which fields are nullable.
type Brother struct {
	PacificId   string `json:"pacificId" validate:"required"` // Required
	FirstName   string `json:"firstName" validate:"required"` // Required
	LastName    string `json:"lastName" validate:"required"`  // Required
	Status      string `json:"status" validate:"required"`    // Required
	Class       string `json:"className"`
	RollCall    string `json:"rollCall"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	BadStanding int    `json:"badStanding"`
}

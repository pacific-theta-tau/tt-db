package models

// TODO: decide which fields are nullable.
type Brother struct {
	RollCall    int    `json:"rollCall"` // Primary Key
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Status      string `json:"status" validate:"required"`
	Class       string `json:"className"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	BadStanding int    `json:"badStanding"`
}

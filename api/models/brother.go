package models

// TODO: decide which fields are nullable.
type Brother struct {
    BrotherID   int    `json:"brotherID"`  // Primary Key
    RollCall    int    `json:"rollCall" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Major       string `json:"major" validate:"required"`
	Status      string `json:"status" validate:"required"`
	Class       string `json:"className"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	BadStanding int    `json:"badStanding"`
}

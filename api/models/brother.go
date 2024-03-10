package models

import "database/sql"

// TODO: decide which fields are nullable.
type Brother struct {
	ID     int            `json:id` // should it be 989?
	First  string         `json:first`
	Last   string         `json:last`
	Status string         `json:status`
	Class  sql.NullString `json:class` // should be required
	Email  sql.NullString `json:email` // nullable
	//Roll_call //nullable
	//Phone_number // nullable
	// created_at?
	// updated_at?
}

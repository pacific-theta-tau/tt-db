package models

import (
	"context"
	"database/sql"
	"time"
)

const brothers_table = "brothers"

// TODO: decide which fields are nullable.
type Brother struct {
	PacificId   string `json:"pacificId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Status      string `json:"status"`
	Class       string `json:"className"`
	RollCall    string `json:"rollCall"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	BadStanding int    `json:"badStanding"`
}

func GetAllBrothers(db *sql.DB) ([]*Brother, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM " + brothers_table
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Scan rows from query to create Brother instances
	var brothers []*Brother
	for rows.Next() {
		var brother Brother
		err = rows.Scan(
			&brother.PacificId,
			&brother.FirstName,
			&brother.LastName,
			&brother.Status,
			&brother.Class,
			&brother.RollCall,
			&brother.Email,
			&brother.PhoneNumber,
			&brother.BadStanding,
		)

		if err != nil {
			return nil, err
		}
		brothers = append(brothers, &brother)
	}

	return brothers, nil
}

// TODO: fix function
func GetBrotherByPacificID(db *sql.DB, pacificID string) (*Brother, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := "SELECT * FROM brothers WHERE pacificId = $1"
	row, err := db.QueryContext(ctx, query, pacificID)
	if err != nil {
		return nil, err
	}

	var brother Brother
	err = row.Scan(
		&brother.PacificId,
		&brother.FirstName,
		&brother.LastName,
		&brother.Status,
		&brother.Class,
		&brother.RollCall,
		&brother.Email,
		&brother.PhoneNumber,
		&brother.BadStanding,
	)

	if err != nil {
		return nil, err
	}

	return &brother, nil
}

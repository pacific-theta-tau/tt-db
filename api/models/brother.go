package models

import (
	"context"
	"database/sql"
	"time"
)

// TODO: decide which fields are nullable.
type Brother struct {
	ID     int    `json:id` // should it be 989?
	First  string `json:first`
	Last   string `json:last`
	Status string `json:status`
	Class  string `json:class` // should be required
	Email  string `json:email` // nullable
	//Roll_call //nullable
	//Phone_number // nullable
	// created_at?
	// updated_at?
}

func GetAllBrothers(db *sql.DB) ([]*Brother, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM brothers"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Scan rows from query to create Brother instances
	var brothers []*Brother
	for rows.Next() {
		var brother Brother
		err = rows.Scan(
			&brother.ID,
			&brother.First,
			&brother.Last,
			&brother.Status,
			&brother.Class,
			&brother.Email,
		)

		if err != nil {
			return nil, err
		}
		brothers = append(brothers, &brother)
	}

	return brothers, nil
}

func GetBrotherByID(db *sql.DB, id string) (*Brother, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel()

	query := "SELECT * FROM brothers WHERE id = @ID"
	row, err := db.Query(query, sql.Named("ID", id))
	if err != nil {
		return nil, err
	}

	var brother *Brother
	err = row.Scan(
		&brother.ID,
		&brother.First,
		&brother.Last,
		&brother.Status,
		&brother.Class,
		&brother.Email,
	)
	if err != nil {
		return nil, err
	}

	return brother, nil
}

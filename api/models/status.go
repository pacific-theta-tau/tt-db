package models

import (
	"database/sql"
)

type Status struct {
    Semester string `json:"semesterStatus"`
    Status string `json:"status"`
}

// Helper function to scan SQL row and create new Brother instance
func CreateStatusFromRow(row *sql.Rows) (Status, error) {
    var status Status
	err := row.Scan(
        &status.Semester,
		&status.Status,
	)
	if err != nil {
		return Status{}, err
	}

	return status, err
}

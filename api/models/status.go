package models

import (
	"database/sql"
)

type Status struct {
    Semester string `json:"semesterLabel"`
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

type BrotherStatus struct {
    BrotherID   int `json:"brotherID"`
    RollCall    string `json:"rollCall"`
    FirstName   string `json:"firstName"`
    LastName    string `jsong:"lastName"`
    Status      string `json:"status"`
    Semester    string `json:"semesterLabel"`
}

func CreateBrotherStatusFromRow(row *sql.Rows) (BrotherStatus, error) {
    var brotherStatus BrotherStatus
    err := row.Scan(
        &brotherStatus.BrotherID,
        &brotherStatus.RollCall,
        &brotherStatus.FirstName,
        &brotherStatus.LastName,
        &brotherStatus.Status,
        &brotherStatus.Semester,
    )
    if err != nil {
        return BrotherStatus{}, err
    }

    return brotherStatus, err
}

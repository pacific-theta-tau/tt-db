package models

import "database/sql"

type Semester struct {
    SemesterID      string `json:"semesterID"`
    SemesterLabel   string `json:"semesterLabel"`
}

// Helper function to scan SQL row and create new Brother instance
func CreateSemesterFromRow(row *sql.Rows) (Semester, error) {
    var semester Semester
	err := row.Scan(
        &semester.SemesterID,
		&semester.SemesterLabel,
	)
	if err != nil {
		return Semester{}, err
	}

	return semester, err
}

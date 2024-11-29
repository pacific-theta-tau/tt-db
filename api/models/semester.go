package models

import "database/sql"

type Semester struct {
    SemesterID      string `json:"semesterID"`
    SemesterLabel   string `json:"semesterLabel"`
}

// Helper function to scan SQL row and create new Semester instance
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

type BrotherStatusFromSemester struct {
    BrotherID   int    `json:"brotherID"`
    RollCall    string `json:"rollCall"`
    FirstName   string `json:"firstName"`
    LastName    string `json:"lastName"`
    Major       string `json:"major"`
    ClassName   string `json:"class"`
    Status      string `json:"status"`
    SemesterID  int `json:"semesterID"`
    SemesterLabel string `json:"semesterLabel"`
}

// Helper function to scan SQL row and create new BrotherStatusFromSemester instance 
func CreateBrotherStatusFromSemesterFromRow(row *sql.Rows) (BrotherStatusFromSemester, error) {
    var b BrotherStatusFromSemester
	err := row.Scan(
        &b.BrotherID,
        &b.RollCall,
        &b.FirstName,
        &b.LastName,
        &b.Major,
        &b.ClassName,
        &b.Status,
        &b.SemesterID,
        &b.SemesterLabel,
	)
	if err != nil {
		return BrotherStatusFromSemester{}, err
	}

	return b, err
}

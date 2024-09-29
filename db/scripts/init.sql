-- initialize tables for dev postgres db
CREATE TYPE status AS ENUM ('Active', 'Pre-Alumnus', 'Alumnus', 'Co-op', 'Transferred', 'Expelled');

CREATE TABLE IF NOT EXISTS brothers(
    brotherID SERIAL PRIMARY KEY,
    rollCall INTEGER NOT NULL,
    firstName TEXT NOT NULL, 
    lastName TEXT NOT NULL, 
    status status NOT NULL, 
    className TEXT DEFAULT '', 
    email TEXT DEFAULT '',
    phoneNumber TEXT DEFAULT 0, 
    badStanding INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS eventsCategory(
    categoryID SERIAL PRIMARY KEY, 
    categoryName TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS events(
    eventID SERIAL PRIMARY KEY, 
    eventName TEXT NOT NULL,
    categoryID INT REFERENCES eventsCategory(categoryID),
    eventLocation Text NOT NULL,
    eventDate date NOT NULL
);

CREATE TABLE IF NOT EXISTS attendance(
	brotherID INT REFERENCES brothers(brotherID),
	eventID INT REFERENCES events(eventID),
	attendanceStatus VARCHAR(20) CHECK (attendanceStatus IN ('present', 'absent', 'excused')),
	PRIMARY KEY (brotherID, eventID)  -- compound PK
);

CREATE TABLE IF NOT EXISTS semester(
    semesterID SERIAL PRIMARY KEY,
    semesterLabel VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS brotherStatus(
    brotherID INT REFERENCES brothers(brotherID),
    semesterID INT REFERENCES semester(semesterID),
    status status,
    PRIMARY KEY (brotherID, semesterID)
);


-- mock entries for testing
INSERT INTO brothers (rollCall, firstName, lastName, status, className, email, phoneNumber, badStanding)
VALUES
    (1, 'John', 'Doe', 'Active', 'Chi', 'na@gmail.com', '(123) 456-7890', 0)
;

INSERT INTO eventsCategory (categoryID, categoryName)
VALUES
    (1, 'Professional Development'),
    (2, 'Brotherhood'),
    (3, 'Community Service')
;

INSERT INTO events (eventName, categoryID, eventLocation, eventDate)
VALUES 
    ('CO-OP Panel', 1, 'Regent Room', '1/28/24'),
    ('Movies', 2, 'CTC', '3/14/24')
;

INSERT INTO semester (semesterLabel) VALUES ('Spring 2023'), ('Fall 2023'), ('Spring 2024'), ('Fall 2024');

INSERT INTO brotherStatus (brotherID, semesterID, status)
VALUES
    (1, 1, 'Active'),
    (1, 2, 'Co-op'),
    (1, 3, 'Active'),
    (1, 4, 'Alumni')
;

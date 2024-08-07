-- initialize tables for dev postgres db
CREATE TYPE status AS ENUM ('Active', 'Pre-Alumnus', 'Alumnus', 'Co-op', 'Transferred', 'Expelled');

CREATE TABLE IF NOT EXISTS brothers(
  rollCall INTEGER PRIMARY KEY NOT NULL,
  firstName TEXT NOT NULL, 
  lastName TEXT NOT NULL, 
  status status NOT NULL, 
  className TEXT DEFAULT '', 
  email TEXT DEFAULT '',
  phoneNumber TEXT DEFAULT 0, 
  badStanding INT DEFAULT 0
);

INSERT INTO brothers (rollCall, firstName, lastName, status, className, email, phoneNumber, badStanding)
VALUES (239, 'Nicolas', 'Ahn', 'Active', 'Chi', 'na@gmail.com', '(123) 456-7890', 0);


CREATE TABLE IF NOT EXISTS eventsCategory(
  categoryID INT PRIMARY KEY NOT NULL, 
  categoryName TEXT NOT NULL
);

INSERT INTO eventsCategory (categoryID, categoryName)
VALUES
    (1, 'Professional Development'),
    (2, 'Brotherhood'),
    (3, 'Community Service'); 


CREATE TABLE IF NOT EXISTS events(
  eventID INT PRIMARY KEY NOT NULL, 
  eventName TEXT NOT NULL,
  categoryID INT REFERENCES eventsCategory(categoryID),
  eventLocation Text NOT NULL,
  eventDate date NOT NULL
);

INSERT INTO events (eventID, eventName, categoryID, eventLocation, eventDate)
VALUES 
    (1, 'CO-OP Panel', 1, 'Regent Room', '1/28/24'),
    (2, 'Movies', 2, 'CTC', '3/14/24');

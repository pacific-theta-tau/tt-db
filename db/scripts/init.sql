-- initialize tables for dev postgres db
CREATE TYPE status AS ENUM ('Active', 'Pre-Alumni', 'Alumni');

CREATE TABLE IF NOT EXISTS brothers(
  pacificId INT PRIMARY KEY NOT NULL, 
  firstName TEXT NOT NULL, 
  lastName TEXT NOT NULL, 
  status status NOT NULL, 
  className TEXT DEFAULT '', 
  rollCall TEXT DEFAULT '',
  email TEXT DEFAULT '',
  phoneNumber TEXT DEFAULT 0, 
  badStanding INT DEFAULT 0
);

INSERT INTO brothers (pacificId, firstName, lastName, status, className, rollCall, email, phoneNumber, badStanding)
VALUES (989123456, 'Nicolas', 'Ahn', 'Active', 'Chi', '000', 'na@gmail.com', '(123) 456-7890', 0);

CREATE TABLE Events(
  event_id INT PRIMARY KEY NOT NULL, 
  event_name TEXT NOT NULL,
  category INT NOT NULL,
  locations Text NOT NULL,
  dates date NOT NULL
);

INSERT INTO Events (event_id, event_name, category, locations, dates)
VALUES (1, 'Meet the Bros', 5, 'Ballroom', '1/28/24');


CREATE TABLE eventsCategory(
  categoryID INT PRIMARY KEY NOT NULL, 
  categoryName TEXT NOT NULL
);

INSERT INTO eventsCategory (categoryID, categoryName)
VALUE (1, 'CO-OP Panel'); 
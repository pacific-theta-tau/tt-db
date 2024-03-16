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
VALUES 
  (989123456, 'Nicolas', 'Ahn', 'Active', 'Chi', '000', 'na@gmail.com', '(123) 456-7890', 0)
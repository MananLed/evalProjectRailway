CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY, 
    firstname TEXT NOT NULL,
    middlename TEXT,
    lastname TEXT NOT NULL, 
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    mobile TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'passenger'))
);

CREATE TABLE IF NOT EXISTS trains (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    destination TEXT NOT NULL,
    departure TIMESTAMP NOT NULL,
    arrival TIMESTAMP NOT NULL,
    totalseats INT NOT NULL,
    bookedseats INT NOT NULL DEFAULT 0,
    seatfare INT NOT NULL 
);

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY,
    trainid UUID NOT NULL,
    passengerid UUID NOT NULL,
    source TEXT NOT NULL,
    destination TEXT NOT NULL,
    departure TIMESTAMP NOT NULL,
    totalfare INT NOT NULL,
    bookedseats INT NOT NULL,

    CONSTRAINT fk_train FOREIGN KEY (trainid) REFERENCES trains(id) ON DELETE CASCADE,
    CONSTRAINT fk_passenger FOREIGN KEY (passengerid) REFERENCES users(id) ON DELETE CASCADE
);
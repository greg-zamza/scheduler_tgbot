CREATE DATABASE app_db;
\c app_db;

CREATE TABLE units (
    name TEXT PRIMARY KEY
);

CREATE TABLE users (
    id BIGINT primary key,
    unit TEXT REFERENCES units(name) ON DELETE CASCADE,
    is_admin INT
);

CREATE TABLE classes (
    day DATE,
    num INT,
    unit TEXT REFERENCES units(name) ON DELETE CASCADE,
    name TEXT,
    room TEXT,
    PRIMARY KEY (day, num, unit)
);

-- insert test values
-- INSERT INTO...
INSERT INTO units(name) VALUES(''); -- for admins;
INSERT INTO users(id, unit, is_admin) VALUES(6271467096, '', 1);


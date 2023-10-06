/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE Users (
	id serial PRIMARY KEY,
	full_name VARCHAR ( 60 ) NOT NULL,
  phone_number VARCHAR ( 15 ) NOT NULL,
  password CHAR ( 60 ) NOT NULL,
  successful_login INT NOT NULL
);

INSERT INTO Users (full_name, phone_number, password, successful_login) VALUES ('Alpen Halim', '+6285883949378', '$2a$12$n9p4U0CmPdCQHJsiP/Xlf.lj5MV6fmqh2nnOnGj6XmEswkaAxnayO', 0)

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE peoples (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    surname varchar(100) NOT NULL,
    patronymic varchar(100),
    age integer NOT NULL,
    gender gender NOT NULL
);

CREATE TABLE person_countries (
    id serial PRIMARY KEY,
    person_id integer NOT NULL,
    country_id varchar(2) NOT NULL,
    probability numeric(5,4),
    FOREIGN KEY (person_id) REFERENCES peoples(id) ON DELETE CASCADE,
    CONSTRAINT valid_probability CHECK (probability IS NULL OR (probability >= 0 AND probability <= 1))
);
CREATE TABLE IF NOT EXISTS people (
           id serial PRIMARY KEY,
           name VARCHAR(255) NOT NULL,
           surname VARCHAR(255) NOT NULL,
           patronymic VARCHAR(255),
           age INT,
           gender VARCHAR(10),
           nationality VARCHAR(255),
           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);



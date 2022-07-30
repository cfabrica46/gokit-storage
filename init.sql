DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password  VARCHAR(128) NOT NULL,
    email VARCHAR(64) NOT NULL UNIQUE
);

INSERT INTO users(username, password,email)
    VALUES
        ('cesar',	'c565fe03ca9b6242e01dfddefe9bba3d98b270e19cd02fd85ceaf75e2b25bf12',	'cesar@gmail.com'),
        ('luis',	'5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5',	'luis@gmail.com')

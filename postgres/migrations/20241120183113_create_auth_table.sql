-- +goose Up
-- +goose StatementBegin

CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(150) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role_id INT NOT NULL,
                       FOREIGN KEY (role_id) REFERENCES roles (id)
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table roles;
-- +goose StatementEnd

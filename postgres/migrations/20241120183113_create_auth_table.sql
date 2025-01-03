-- +goose Up
-- +goose StatementBegin

CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       role_name VARCHAR(50) NOT NULL
);

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(150) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) ,
                       role_id INT DEFAULT 1 NOT NULL,
                       FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE SET DEFAULT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NULL -- write a trigger to update time automatically
);


INSERT INTO roles (role_name) VALUES ('admin'), ('user');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table roles;
-- +goose StatementEnd

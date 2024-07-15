-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login varchar(100) UNIQUE NOT NULL, 
			password varchar(200) NOT NULL
		);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			user_id int NOT NULL,
			secret bytea UNIQUE NOT NULL, 
			meta varchar NOT NULL,
			comment bytea			
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			user_id int NOT NULL,
			secret VARCHAR(8000) UNIQUE NOT NULL, 
			meta varchar NOT NULL,
			comment varchar(8000)			
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
-- +goose StatementEnd

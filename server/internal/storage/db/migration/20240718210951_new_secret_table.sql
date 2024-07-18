-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			user_id int NOT NULL,
			secret varchar UNIQUE NOT NULL, 
			meta varchar NOT NULL,
			comment varchar			
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login varchar(100) UNIQUE NOT NULL, 
			password varchar(200) NOT NULL
		);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);

CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			user_id int NOT NULL,
			secret VARCHAR(8000) UNIQUE NOT NULL, 
			meta varchar NOT NULL,
			comment varchar(8000)			
		);
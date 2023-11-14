-- BEGIN;

--Set timezone
-- SET TIME ZONE 'Asia/Bangkok';

-- CREATE OR REPLACE FUNCTION set_updated_at_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = now();
--     RETURN NEW;   
-- END;
-- $$ language 'plpgsql';

CREATE TABLE users  (
  id VARCHAR(7) PRIMARY KEY NOT NULL,
  username VARCHAR UNIQUE NOT NULL,
  password VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- END;
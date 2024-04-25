ALTER TABLE users
    DROP COLUMN connected_state;

ALTER TABLE users
    CHANGE COLUMN email username VARCHAR(255) NOT NULL UNIQUE;

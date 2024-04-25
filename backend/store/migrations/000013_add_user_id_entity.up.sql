ALTER TABLE entity
    ADD COLUMN user_id INT,
    ADD CONSTRAINT fk_entity_user_id FOREIGN KEY (user_id) REFERENCES users(id);
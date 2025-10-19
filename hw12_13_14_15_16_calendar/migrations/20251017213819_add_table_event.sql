-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    event_time TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL,
    description VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    time_to_notify TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd

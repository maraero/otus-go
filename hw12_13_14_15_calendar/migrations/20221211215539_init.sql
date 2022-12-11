-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id INTEGER AUTO_INCREMENT,
    title TEXT NOT NULL,
    date_start TIMESTAMP NOT NULL,
    date_end TIMESTAMP NOT NULL,
    description TEXT NULL,
    user_id STRING NUT NULL,
    date_notification TIMESTAMP NULL,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd

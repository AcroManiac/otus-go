-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id UUID primary key,
    title text not null,
    description text,
    owner text not null,
    start_time timestamptz not null,
    duration interval not null,
    notify interval
);

CREATE INDEX start_time_owner_idx ON events USING btree (start_time, owner);

GRANT ALL PRIVILEGES ON TABLE events IN SCHEMA public TO dbuser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd

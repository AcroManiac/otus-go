-- +goose Up
-- +goose StatementBegin
CREATE TABLE notices (
    id UUID primary key,
    send_time timestamptz not null
);

CREATE INDEX notices_id_idx ON notices USING btree (id);

GRANT ALL PRIVILEGES ON TABLE notices IN SCHEMA public TO dbuser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notices;
-- +goose StatementEnd

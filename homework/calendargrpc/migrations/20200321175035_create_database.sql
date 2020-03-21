-- +goose Up
-- +goose StatementBegin
CREATE DATABASE calendar OWNER dbuser;
GRANT ALL PRIVILEGES ON DATABASE calendar TO dbuser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE calendar;
-- +goose StatementEnd

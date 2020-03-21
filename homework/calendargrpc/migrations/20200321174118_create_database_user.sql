-- +goose Up
-- +goose StatementBegin
CREATE USER dbuser WITH encrypted password 'En9NR2b869';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER dbuser;
-- +goose StatementEnd

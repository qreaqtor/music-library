-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table songs
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    group varchar(100) NOT NULL,
    song varchar(100) NOT NULL,
    releaseDate DATE DEFAULT current_date,
    link text,
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

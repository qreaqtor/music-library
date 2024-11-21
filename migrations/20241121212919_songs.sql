-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table songs
(
    group varchar(100) NOT NULL,
    song varchar(100) NOT NULL
    releaseDate DATE DEFAULT current_date,
    link text,
    lyrics text[],
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

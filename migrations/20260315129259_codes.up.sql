create table codes
(
    id              uuid primary key,
    code            varchar(20) not null,
    percent         int         not null,
    amount_of_usage int,
    expired_at      timestamp,
    created_at      timestamp default current_timestamp,
    deleted_at      timestamp default null
);

CREATE UNIQUE INDEX codes_code_unique_active
    ON codes (code)
    WHERE deleted_at IS NULL;

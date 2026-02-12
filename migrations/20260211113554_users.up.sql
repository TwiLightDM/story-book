create table users
(
    id         uuid primary key,
    name       varchar(50)  not null,
    surname    varchar(50)  not null,
    email      varchar(50)  not null,
    phone      varchar(11)  not null,
    password   varchar(100) not null,
    salt       varchar(20)  not null,
    role       varchar(20)  not null,
    question   text         not null,
    answer     text         not null,
    points     int       default 0,
    created_at timestamp default current_timestamp,
    deleted_at timestamp default null
);

CREATE UNIQUE INDEX users_email_unique_active
    ON users (email)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX users_phone_unique_active
    ON users (phone)
    WHERE deleted_at IS NULL;
create table books
(
    id          uuid primary key,
    title       varchar(100)   not null,
    author      varchar(100)   not null,
    year        smallint       not null,
    cost        numeric(10, 2) not null,
    discount    smallint,
    publisher   varchar(100)   not null,
    description text,
    amount      int       default 0,
    image_data  bytea,
    image_mime  varchar(20),
    created_at  timestamp default current_timestamp,
    deleted_at  timestamp default null
)
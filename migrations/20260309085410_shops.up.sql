create table shops
(
    id      uuid primary key,
    name    text        not null,
    address text unique not null
)
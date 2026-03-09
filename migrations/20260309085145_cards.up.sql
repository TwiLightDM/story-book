create table cards
(
    id              uuid primary key,
    number_of_card  varchar(16) unique         not null,
    expiration_date timestamp                  not null,
    cvv             varchar(3)                 not null,
    user_id         uuid references users (id) not null
)

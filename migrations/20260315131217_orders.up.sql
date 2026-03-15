create table orders
(
    id            uuid primary key,
    status        varchar(10) check (status in ('created', 'paid', 'received')) default 'created',
    delivery_type varchar(10) check (delivery_type in ('delivery', 'pickup')),
    address       text,
    cost          numeric(10, 2)             not null,
    points        int                        not null,
    user_id       uuid references users (id) not null,
    shop_id       uuid references shops (id),
    code_id       uuid references codes (id),
    created_at    timestamp                                                     default current_timestamp,
    deleted_at    timestamp                                                     default null
)
create table order_books
(
    id            uuid primary key,
    amount        int                         not null,
    price_for_one numeric(10, 2)              not null,
    book_id       uuid references books (id)  not null,
    order_id      uuid references orders (id) not null
)
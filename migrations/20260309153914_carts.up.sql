create table carts
(
    id      uuid primary key,
    amount  int default 0,
    user_id uuid references users (id) not null,
    book_id uuid references books (id) not null,
    UNIQUE (user_id, book_id)
)


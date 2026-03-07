create table favourites
(
    id      uuid primary key,
    book_id uuid references books (id),
    user_id uuid references users (id),
    created_at  timestamp default current_timestamp,
    deleted_at  timestamp default null
)
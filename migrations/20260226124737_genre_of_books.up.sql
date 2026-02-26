create table genre_of_books(
    id uuid primary key,
    genre varchar(30) not null,
    book_id uuid references books(id) not null,
    created_at timestamp default current_timestamp,
    deleted_at timestamp default null
)
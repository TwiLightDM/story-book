create table ratings
(
    id      uuid primary key,
    stars   int CHECK (stars BETWEEN 1 AND 5) not null,
    user_id uuid references users (id)        not null,
    book_id uuid references books (id)        not null
)
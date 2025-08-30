-- +goose Up
create table posts (
  id uuid primary key,
  created_at timestamp not null default now(),
  updated_at timestamp,
  title text not null,
  url text unique not null,
  description text not null,
  published_at timestamp,
  feed_id uuid not null,
  foreign key (feed_id) references feeds(id)
);

-- +goose Down
drop table posts;

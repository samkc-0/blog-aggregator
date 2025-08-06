-- +goose Up
create table users (
  id bigserial primary key,
  name varchar(64) not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

-- +goose Down
drop table users;

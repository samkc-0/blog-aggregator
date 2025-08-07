-- name: CreateUser :one
insert into users (
  id, name, created_at, updated_at
) values (
  $1,
  $2,
  $3,
  $4
) returning *;


-- name: GetUser :one
select * from users
where name = $1 limit 1;

-- name: GetUsers :many
select * from users;

-- name: DeleteAllUsers :exec
delete from users;

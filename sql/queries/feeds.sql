-- name: CreateFeed :one
insert into feeds (
  id, created_at, updated_at, name, url, user_id
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) returning *;


-- name: GetFeed :one
select * from feeds
where id = $1 limit 1;

-- name: GetFeedByUrl :one
select * from feeds
where url = $1 limit 1;

-- name: GetFeeds :many
select * from feeds;

-- name: MarkFeedFetched :exec
update feeds
set
last_fetched_at = now(),
updated_at = now()
where id = $1;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at nulls first
limit 1;

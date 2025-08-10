-- name: CreateFeedFollow :one
with follows as (
  insert into feed_follows (created_at, updated_at, user_id, feed_id)
  values ($1, $2, $3, $4)
  returning *
) select
  follows.*,
  users.name as user_name,
  feeds.name as feed_name
from follows
join users on follows.user_id = users.id
join feeds on follows.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
select *, feeds.name as name from feed_follows
join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1;

-- name: CreatePost :one
insert into posts (
  title, url, description, published_at, feed_id
) values (
  $1, $2, $3, $4, $5
) returning *;

-- name: GetPostsForUser :many
select * from posts
inner join feed_follows on
  posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2;

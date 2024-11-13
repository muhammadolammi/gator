-- name: CreatePost :one
INSERT INTO posts ( created_at, updated_at, title,url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT  posts.* 
FROM posts
JOIN feeds ON feeds.id= posts.feed_id
JOIN users ON users.id=feeds.user_id 
WHERE users.id = $1
ORDER BY published_at DESC
LIMIT $2; 

-- name: PostExists :one
SELECT EXISTS (
    SELECT 1
    FROM posts
    WHERE url = $1
) AS exists;

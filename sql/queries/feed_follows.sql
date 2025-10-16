-- name: CreateFeedFollow :many
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
)
SELECT *
FROM inserted id
INNER JOIN users
    ON i.user_id = users.id
INNER JOIN feeds
    ON i.feed_id = feeds.i;
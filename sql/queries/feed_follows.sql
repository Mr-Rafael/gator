-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING id, user_id, feed_id
)
SELECT *
FROM inserted i
INNER JOIN users
    ON i.user_id = users.id
INNER JOIN feeds
    ON i.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.user_id, feed_follows.feed_id, name, url
FROM feed_follows
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;
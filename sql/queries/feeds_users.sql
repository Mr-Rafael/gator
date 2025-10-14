-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;
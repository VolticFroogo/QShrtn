-- name: redirect-from-id
SELECT id, url FROM redirect WHERE id=?;

-- name: redirect-from-url
SELECT id, url FROM redirect WHERE url=?;

-- name: insert-redirect
INSERT INTO redirect (id, url) VALUES (?, ?);
-- name: redirect-from-id
SELECT id, url FROM redirect WHERE BINARY id=?;

-- name: redirect-from-url
SELECT id, url FROM redirect WHERE BINARY url=?;

-- name: insert-redirect
INSERT INTO redirect (id, url) VALUES (?, ?);
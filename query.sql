-- name: ListThings :many
SELECT * FROM Things
ORDER BY Year DESC, Month DESC, Day DESC;

-- name: InsertThing :execlastid
INSERT INTO Things(Year, Month, Day, What) VALUES(?, ?, ?, ?);

-- name: DeleteThing :exec
DELETE FROM Things WHERE id=?;

-- name: ListTopThings :many
SELECT COUNT(*) as cnt, What FROM Things
GROUP BY What
ORDER BY cnt DESC
LIMIT 10;
-- name: CreateEntry :one
INSERT INTO entires (
    account_id, amount
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entires
WHERE id = $1 LIMIT 1;

-- name: ListEntry :many
SELECT * FROM entires
         WHERE account_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3;


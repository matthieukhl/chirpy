-- name: GetChirpByID :one
SELECT 
    id, body, user_id, created_at, updated_at
FROM 
    chirps
WHERE 
    id = $1;
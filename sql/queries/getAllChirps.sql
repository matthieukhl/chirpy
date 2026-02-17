-- name: GetAllChirps :many
SELECT 
    id, body, user_id, created_at, updated_at
FROM 
    chirps
ORDER BY 
    created_at ASC;
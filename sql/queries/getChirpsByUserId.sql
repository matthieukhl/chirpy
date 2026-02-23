-- name: GetChirpsByUserId :many
SELECT 
    id, body, user_id, created_at, updated_at
FROM 
    chirps
WHERE 
    user_id = $1
ORDER BY 
    created_at ASC;

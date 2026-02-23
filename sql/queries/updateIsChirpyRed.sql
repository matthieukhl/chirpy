-- name: UpgradeChirpyRed :one
UPDATE 
    users 
SET 
    is_chirpy_red = true,
    updated_at = NOW()
WHERE 
    id = $1
RETURNING 
    id, created_at, updated_at, email, is_chirpy_red;
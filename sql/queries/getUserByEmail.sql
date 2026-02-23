-- name: GetUserByEmail :one 
SELECT 
    id, 
    email,
    hashed_password,
    is_chirpy_red,
    created_at,
    updated_at
FROM 
    users
WHERE 
    email = $1;
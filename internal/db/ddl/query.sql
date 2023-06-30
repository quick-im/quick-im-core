-- name: GetConvercationSessionCountByConvercationPkId :one
SELECT count(id) FROM conversation_session_id
WHERE convercation_id = @convercation_id::uuid;
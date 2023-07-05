-- name: GetConvercationSessionCountByConvercationPkId :one
SELECT count(id) FROM public.conversation_session_id
WHERE convercation_id = @convercation_id::varchar;

-- name: CreateConvercation :exec
INSERT INTO public.conversations
(conversation_id, last_msg_id, last_send_time, is_delete, is_archive, conversation_type, last_send_session)
VALUES(@convercation_id::varchar, NULL, NULL, false, false, 0, NULL);

-- name: SessionJoinsConvercationUseCopyFrom :copyfrom
INSERT INTO public.conversation_session_id
(session_id, convercation_id)
VALUES($1, $2);
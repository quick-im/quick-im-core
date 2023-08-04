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

-- name: GetJoinedConversations :many
SELECT id, session_id, last_recv_msg_id, is_kick_out, convercation_id, 0 as unread
FROM public.conversation_session_id WHERE session_id = @session_id::varchar AND is_kick_out = false;

-- name: GetJoinedConversationsUnReadMsgCount :one
SELECT count(msg_id) as unread
FROM public.messages WHERE msg_id BETWEEN @last_recv_msg_id::varchar AND @last_send_msg_id::varchar;

-- name: GetConversationUnReadMsgCount :one
SELECT count(msg_id) as unread
FROM public.messages WHERE convercation_id = @convercation_id::varchar AND msg_id BETWEEN @last_recv_msg_id::varchar AND @last_send_msg_id::varchar;

-- name: GetConversationInfo :one
SELECT conversation_id, last_msg_id, last_send_time, is_delete, conversation_type, last_send_session, is_archive
FROM public.conversations WHERE conversation_id = @conversation_id::varchar;

-- name: CheckJoinedonversation :one
SELECT count(id)
FROM public.conversation_session_id WHERE session_id = @session_id::varchar AND convercation_id= @convercation_id::varchar AND is_kick_out = false;

-- name: DeleteConversations :batchexec
UPDATE public.conversations
SET is_delete=true
WHERE conversation_id=$1;

-- name: ArchiveConversations :batchexec
UPDATE public.conversations
SET is_archive=true
WHERE conversation_id=$1;

-- name: UnArchiveConversations :batchexec
UPDATE public.conversations
SET is_archive=false
WHERE conversation_id=$1;

-- name: KickoutForConversation :batchexec
UPDATE public.conversation_session_id
SET is_kick_out=true
WHERE session_id = $1 AND convercation_id=$2;

-- name: UpdateConversationLastMsg :exec
UPDATE public.conversations
SET last_msg_id= @last_msg_id::varchar, last_send_time=$1, last_send_session= @last_send_session::varchar
WHERE conversation_id= @conversation_id::varchar;

-- name: SaveMsgToDb :batchexec
INSERT INTO public.messages
(msg_id, convercation_id, from_session, send_time, status, "type", "content")
VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: GetMsgFromDbInRange :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = @convercation_id::text BETWEEN @start_msg_id::text AND @end_msg_id::text;

-- name: GetLast30MsgFromDb :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = @convercation_id::text ORDER BY msg_id DESC LIMIT 30;

-- name: GetThe30MsgBeforeTheId :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = @convercation_id::text AND msg_id < @msg_id::text ORDER BY msg_id DESC LIMIT 30;

-- name: GetThe30MsgAfterTheId :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = @convercation_id::text AND msg_id > @msg_id::text ORDER BY msg_id ASC LIMIT 30;
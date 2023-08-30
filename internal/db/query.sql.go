// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: query.sql

package db

import (
	"context"
	"time"
)

const checkJoinedonversation = `-- name: CheckJoinedonversation :one
SELECT count(id)
FROM public.conversation_session_id WHERE session_id = $1::varchar AND conversation_id= $2::varchar AND is_kick_out = false
`

type CheckJoinedonversationParams struct {
	SessionID      string
	ConversationID string
}

func (q *Queries) CheckJoinedonversation(ctx context.Context, arg CheckJoinedonversationParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkJoinedonversation, arg.SessionID, arg.ConversationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createConversation = `-- name: CreateConversation :exec
INSERT INTO public.conversations
(conversation_id, last_msg_id, last_send_time, is_delete, is_archive, conversation_type, last_send_session)
VALUES($1::varchar, NULL, NULL, false, false, 0, NULL)
`

func (q *Queries) CreateConversation(ctx context.Context, conversationID string) error {
	_, err := q.db.Exec(ctx, createConversation, conversationID)
	return err
}

const getConversationInfo = `-- name: GetConversationInfo :one
SELECT conversation_id, last_msg_id, last_send_time, is_delete, conversation_type, last_send_session, is_archive
FROM public.conversations WHERE conversation_id = $1::varchar
`

func (q *Queries) GetConversationInfo(ctx context.Context, conversationID string) (Conversation, error) {
	row := q.db.QueryRow(ctx, getConversationInfo, conversationID)
	var i Conversation
	err := row.Scan(
		&i.ConversationID,
		&i.LastMsgID,
		&i.LastSendTime,
		&i.IsDelete,
		&i.ConversationType,
		&i.LastSendSession,
		&i.IsArchive,
	)
	return i, err
}

const getConversationSessionCountByConversationPkId = `-- name: GetConversationSessionCountByConversationPkId :one
SELECT count(id) FROM public.conversation_session_id
WHERE conversation_id = $1::varchar
`

func (q *Queries) GetConversationSessionCountByConversationPkId(ctx context.Context, conversationID string) (int64, error) {
	row := q.db.QueryRow(ctx, getConversationSessionCountByConversationPkId, conversationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getConversationUnReadMsgCount = `-- name: GetConversationUnReadMsgCount :one
SELECT count(msg_id) as unread
FROM public.messages WHERE conversation_id = $1::varchar AND msg_id BETWEEN $2::varchar AND $3::varchar
`

type GetConversationUnReadMsgCountParams struct {
	ConversationID string
	LastRecvMsgID  string
	LastSendMsgID  string
}

func (q *Queries) GetConversationUnReadMsgCount(ctx context.Context, arg GetConversationUnReadMsgCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getConversationUnReadMsgCount, arg.ConversationID, arg.LastRecvMsgID, arg.LastSendMsgID)
	var unread int64
	err := row.Scan(&unread)
	return unread, err
}

const getConversationsAllUsers = `-- name: GetConversationsAllUsers :many
SELECT id, session_id, last_recv_msg_id, is_kick_out, conversation_id
FROM public.conversation_session_id WHERE conversation_id = conversation_id::text
`

func (q *Queries) GetConversationsAllUsers(ctx context.Context) ([]ConversationSessionID, error) {
	rows, err := q.db.Query(ctx, getConversationsAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ConversationSessionID
	for rows.Next() {
		var i ConversationSessionID
		if err := rows.Scan(
			&i.ID,
			&i.SessionID,
			&i.LastRecvMsgID,
			&i.IsKickOut,
			&i.ConversationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJoinedConversations = `-- name: GetJoinedConversations :many
SELECT id, session_id, last_recv_msg_id, is_kick_out, conversation_id, 0 as unread
FROM public.conversation_session_id WHERE session_id = $1::varchar AND is_kick_out = false
`

type GetJoinedConversationsRow struct {
	ID             int32
	SessionID      string
	LastRecvMsgID  *string
	IsKickOut      bool
	ConversationID string
	Unread         int32
}

func (q *Queries) GetJoinedConversations(ctx context.Context, sessionID string) ([]GetJoinedConversationsRow, error) {
	rows, err := q.db.Query(ctx, getJoinedConversations, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetJoinedConversationsRow
	for rows.Next() {
		var i GetJoinedConversationsRow
		if err := rows.Scan(
			&i.ID,
			&i.SessionID,
			&i.LastRecvMsgID,
			&i.IsKickOut,
			&i.ConversationID,
			&i.Unread,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJoinedConversationsUnReadMsgCount = `-- name: GetJoinedConversationsUnReadMsgCount :one
SELECT count(msg_id) as unread
FROM public.messages WHERE msg_id BETWEEN $1::varchar AND $2::varchar
`

type GetJoinedConversationsUnReadMsgCountParams struct {
	LastRecvMsgID string
	LastSendMsgID string
}

func (q *Queries) GetJoinedConversationsUnReadMsgCount(ctx context.Context, arg GetJoinedConversationsUnReadMsgCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getJoinedConversationsUnReadMsgCount, arg.LastRecvMsgID, arg.LastSendMsgID)
	var unread int64
	err := row.Scan(&unread)
	return unread, err
}

const getLast30MsgFromDb = `-- name: GetLast30MsgFromDb :many
SELECT msg_id, conversation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE conversation_id = $1::text ORDER BY msg_id DESC LIMIT 30
`

func (q *Queries) GetLast30MsgFromDb(ctx context.Context, conversationID string) ([]Message, error) {
	rows, err := q.db.Query(ctx, getLast30MsgFromDb, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConversationID,
			&i.FromSession,
			&i.SendTime,
			&i.Status,
			&i.Type,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastOneMsgIdFromDb = `-- name: GetLastOneMsgIdFromDb :one
SELECT last_msg_id FROM public.conversations WHERE conversation_id = $1::text
`

func (q *Queries) GetLastOneMsgIdFromDb(ctx context.Context, conversationID string) (*string, error) {
	row := q.db.QueryRow(ctx, getLastOneMsgIdFromDb, conversationID)
	var last_msg_id *string
	err := row.Scan(&last_msg_id)
	return last_msg_id, err
}

const getMsgFromDbInRange = `-- name: GetMsgFromDbInRange :many
SELECT msg_id, conversation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE conversation_id = $1::text BETWEEN $2::text AND $3::text
`

type GetMsgFromDbInRangeParams struct {
	ConversationID string
	StartMsgID     string
	EndMsgID       string
}

func (q *Queries) GetMsgFromDbInRange(ctx context.Context, arg GetMsgFromDbInRangeParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getMsgFromDbInRange, arg.ConversationID, arg.StartMsgID, arg.EndMsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConversationID,
			&i.FromSession,
			&i.SendTime,
			&i.Status,
			&i.Type,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThe30MsgAfterTheId = `-- name: GetThe30MsgAfterTheId :many
SELECT msg_id, conversation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE conversation_id = $1::text AND msg_id > $2::text ORDER BY msg_id ASC LIMIT 30
`

type GetThe30MsgAfterTheIdParams struct {
	ConversationID string
	MsgID          string
}

func (q *Queries) GetThe30MsgAfterTheId(ctx context.Context, arg GetThe30MsgAfterTheIdParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getThe30MsgAfterTheId, arg.ConversationID, arg.MsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConversationID,
			&i.FromSession,
			&i.SendTime,
			&i.Status,
			&i.Type,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThe30MsgBeforeTheId = `-- name: GetThe30MsgBeforeTheId :many
SELECT msg_id, conversation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE conversation_id = $1::text AND msg_id < $2::text ORDER BY msg_id DESC LIMIT 30
`

type GetThe30MsgBeforeTheIdParams struct {
	ConversationID string
	MsgID          string
}

func (q *Queries) GetThe30MsgBeforeTheId(ctx context.Context, arg GetThe30MsgBeforeTheIdParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getThe30MsgBeforeTheId, arg.ConversationID, arg.MsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConversationID,
			&i.FromSession,
			&i.SendTime,
			&i.Status,
			&i.Type,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type SessionJoinsConversationUseCopyFromParams struct {
	SessionID      string
	ConversationID string
}

const updateConversationLastMsg = `-- name: UpdateConversationLastMsg :exec
UPDATE public.conversations
SET last_msg_id= $2::varchar, last_send_time=$1, last_send_session= $3::varchar
WHERE conversation_id= $4::varchar AND last_msg_id < $2::varchar
`

type UpdateConversationLastMsgParams struct {
	LastSendTime    *time.Time
	LastMsgID       string
	LastSendSession string
	ConversationID  string
}

func (q *Queries) UpdateConversationLastMsg(ctx context.Context, arg UpdateConversationLastMsgParams) error {
	_, err := q.db.Exec(ctx, updateConversationLastMsg,
		arg.LastSendTime,
		arg.LastMsgID,
		arg.LastSendSession,
		arg.ConversationID,
	)
	return err
}

const updateSessionLastRecvMsg = `-- name: UpdateSessionLastRecvMsg :exec
UPDATE public.conversation_session_id
SET last_recv_msg_id= $1::varchar
WHERE conversation_id= $2::varchar AND session_id IN ($3::varchar) AND last_recv_msg_id < $1::varchar
`

type UpdateSessionLastRecvMsgParams struct {
	LastMsgID      string
	ConversationID string
	SessionID      string
}

func (q *Queries) UpdateSessionLastRecvMsg(ctx context.Context, arg UpdateSessionLastRecvMsgParams) error {
	_, err := q.db.Exec(ctx, updateSessionLastRecvMsg, arg.LastMsgID, arg.ConversationID, arg.SessionID)
	return err
}

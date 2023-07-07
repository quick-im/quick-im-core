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
FROM public.conversation_session_id WHERE session_id = $1::varchar AND convercation_id= $2::varchar AND is_kick_out = false
`

type CheckJoinedonversationParams struct {
	SessionID      string
	ConvercationID string
}

func (q *Queries) CheckJoinedonversation(ctx context.Context, arg CheckJoinedonversationParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkJoinedonversation, arg.SessionID, arg.ConvercationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createConvercation = `-- name: CreateConvercation :exec
INSERT INTO public.conversations
(conversation_id, last_msg_id, last_send_time, is_delete, is_archive, conversation_type, last_send_session)
VALUES($1::varchar, NULL, NULL, false, false, 0, NULL)
`

func (q *Queries) CreateConvercation(ctx context.Context, convercationID string) error {
	_, err := q.db.Exec(ctx, createConvercation, convercationID)
	return err
}

const getConvercationSessionCountByConvercationPkId = `-- name: GetConvercationSessionCountByConvercationPkId :one
SELECT count(id) FROM public.conversation_session_id
WHERE convercation_id = $1::varchar
`

func (q *Queries) GetConvercationSessionCountByConvercationPkId(ctx context.Context, convercationID string) (int64, error) {
	row := q.db.QueryRow(ctx, getConvercationSessionCountByConvercationPkId, convercationID)
	var count int64
	err := row.Scan(&count)
	return count, err
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

const getJoinedConversations = `-- name: GetJoinedConversations :many
SELECT id, session_id, last_recv_msg_id, is_kick_out, convercation_id
FROM public.conversation_session_id WHERE session_id = $1::varchar AND is_kick_out = false
`

func (q *Queries) GetJoinedConversations(ctx context.Context, sessionID string) ([]ConversationSessionID, error) {
	rows, err := q.db.Query(ctx, getJoinedConversations, sessionID)
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
			&i.ConvercationID,
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

const getLast30MsgFromDb = `-- name: GetLast30MsgFromDb :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = $1::text ORDER BY msg_id DESC LIMIT 30
`

func (q *Queries) GetLast30MsgFromDb(ctx context.Context, convercationID string) ([]Message, error) {
	rows, err := q.db.Query(ctx, getLast30MsgFromDb, convercationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConvercationID,
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

const getMsgFromDbInRange = `-- name: GetMsgFromDbInRange :many
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = $1::text BETWEEN $2::text AND $3::text
`

type GetMsgFromDbInRangeParams struct {
	ConvercationID string
	StartMsgID     string
	EndMsgID       string
}

func (q *Queries) GetMsgFromDbInRange(ctx context.Context, arg GetMsgFromDbInRangeParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getMsgFromDbInRange, arg.ConvercationID, arg.StartMsgID, arg.EndMsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConvercationID,
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
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = $1::text AND msg_id > $2::text ORDER BY msg_id ASC LIMIT 30
`

type GetThe30MsgAfterTheIdParams struct {
	ConvercationID string
	MsgID          string
}

func (q *Queries) GetThe30MsgAfterTheId(ctx context.Context, arg GetThe30MsgAfterTheIdParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getThe30MsgAfterTheId, arg.ConvercationID, arg.MsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConvercationID,
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
SELECT msg_id, convercation_id, from_session, send_time, status, "type", "content"
FROM public.messages WHERE convercation_id = $1::text AND msg_id < $2::text ORDER BY msg_id DESC LIMIT 30
`

type GetThe30MsgBeforeTheIdParams struct {
	ConvercationID string
	MsgID          string
}

func (q *Queries) GetThe30MsgBeforeTheId(ctx context.Context, arg GetThe30MsgBeforeTheIdParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getThe30MsgBeforeTheId, arg.ConvercationID, arg.MsgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.MsgID,
			&i.ConvercationID,
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

type SessionJoinsConvercationUseCopyFromParams struct {
	SessionID      string
	ConvercationID string
}

const updateConversationLastMsg = `-- name: UpdateConversationLastMsg :exec
UPDATE public.conversations
SET last_msg_id= $2::varchar, last_send_time=$1, last_send_session= $3::varchar
WHERE conversation_id= $4::varchar
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

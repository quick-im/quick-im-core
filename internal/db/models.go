// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Conversation struct {
	ConversationID   string
	LastMsgID        *string
	LastSendTime     *time.Time
	IsDelete         bool
	ConversationType int32
	LastSendSession  *string
	IsArchive        bool
}

type ConversationSessionID struct {
	ID             int32
	SessionID      string
	LastRecvMsgID  *string
	IsKickOut      bool
	ConversationID string
}

type Message struct {
	MsgID          string
	ConversationID string
	FromSession    int32
	SendTime       pgtype.Timestamp
	Status         int32
	Type           int32
	Content        *string
}

type Session struct {
	ID      int32
	Session string
}

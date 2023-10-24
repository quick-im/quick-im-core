package model

import (
	"github.com/quick-im/quick-im-core/internal/jtime"
)

type Msg struct {
	MsgId          string     `rethinkdb:"msg_id" imdb:"pk"`
	ConversationID string     `rethinkdb:"conversation_id" imdb:"index"`
	FromSession    string     `rethinkdb:"from_session" imdb:"index"`
	SendTime       jtime.Time `rethinkdb:"send_time" imdb:"index"`
	Status         int32      `rethinkdb:"status"`
	Type           int32      `rethinkdb:"type"`
	Content        string     `rethinkdb:"content"`
}

package msghub

import "time"

type SendMsgArgs struct {
	MsgId          string
	FromSession    string
	ConversationID string
	MsgType        int32
	Content        []byte
	SendTime       time.Time
}

type SendMsgReply struct {
}

package msghub

import "time"

type SendMsgArgs struct {
	MsgId          string
	FromSession    int32
	ConversationID string
	MsgType        int32
	Content        []byte
	SendTime       time.Time
}

type SendMsgReply struct {
}

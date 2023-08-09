package model

type Msg struct {
	MsgId string `rethinkdb:"msg_id" imdb:"pk"`
}

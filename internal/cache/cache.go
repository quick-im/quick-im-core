package cache

import "github.com/quick-im/quick-im-core/internal/db"

type Cache interface {
	AddConversationSessions(conversation string, sessions []string) error
	DelConversationSession(conversation string, session []string) error
	CleanConversation(conversation string) error
	CountConversationSessions(conversation string) (int64, error)
	IsExistsInConversation(conversation, session string) (bool, error)
	GetConversationSessions(conversation string) ([]string, error)
	KeyExistInCache(key string) (bool, error)
	SyncConversationLastMsgId(conversationId, msgId string) error
	GetConversationLastMsgId(conversationId string) (string, error)
	SetConversationInfo(conversationId string, info db.Conversation) error
	GetConversationInfo(conversationId string) (db.Conversation, error)
	UnSetConversationInfo(conversationId string) error
}

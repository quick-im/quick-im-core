package cache

type Cache interface {
	AddConversationSessions(conversation string, sessions []string) error
	DelConversationSession(conversation string, session []string) error
	CleanConversation(conversation string) error
	CountConversationSessions(conversation string) (int64, error)
	IsExistsInConversation(conversation, session string) (bool, error)
	GetConversationSessions(conversation string) ([]string, error)
	KeyExistInCache(key string) (bool, error)
	SyncConversationLastMsgId(conversationId, msgId string) error
}

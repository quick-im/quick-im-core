package cache

type Cache interface {
	AddConversationSessions(conversation string, sessions []string) error
	DelConversationSession(conversation string, session []string) error
	CleanConversation(conversation string) error
	CountConversationSessions(conversation string)
	IsExistsInConversation(conversation, session string) (bool, error)
	GetConversationSessions(conversation string) ([]string, error)
	KeyExistInCache(key string) (bool, error)
	SyncConversationLastMsgId(conversationId, msgId string) error
	GetConversationLastMsgId(conversationId string) (string, error)
}

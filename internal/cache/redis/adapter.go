package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type adapter struct {
	*redis.Client
}

func (a *adapter) AddConverstaionSessions(conversation string, sessions []string) error {
	return a.SAdd(context.Background(), conversation, sessions).Err()
}

func (a *adapter) DelConversationSession(conversation string, session []string) error {
	return a.SRem(context.Background(), conversation, session).Err()
}

func (a *adapter) CleanConversation(conversation string) error {
	return a.Del(context.Background(), conversation).Err()
}

func (a *adapter) CountConversationSessions(conversation string) (int64, error) {
	return a.SCard(context.Background(), conversation).Result()
}

func (a *adapter) IsExistsInConversation(conversation, session string) (bool, error) {
	return a.SIsMember(context.Background(), conversation, session).Result()
}

func (a *adapter) GetConversationSessions(conversation string) ([]string, error) {
	return a.SMembers(context.Background(), conversation).Result()
}

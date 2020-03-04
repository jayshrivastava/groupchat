package repositories

import(
	chat "github.com/jayshrivastava/groupchat/proto"
)

type ChannelRepository interface {
	Create(key string) error
	Get(key string) (chan chat.StreamResponse, error)
}

type GroupRepository interface {
	GetGroupMembers(group string, username string) ([]string, error)
	CreateIfNotExists(group string) (error)
	AddUserToGroup(username string, group string) (error)
	RemoveUserFromGroup(username string, group string) (error)
}

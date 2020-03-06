package repositories

import (
	chat "github.com/jayshrivastava/groupchat/proto"
)

type ChannelRepository interface {
	Create(key string) error
	Get(key string) (chan chat.StreamResponse, error)
}

type GroupRepository interface {
	GetGroupMembers(group string, username string) ([]string, error)
	CreateIfNotExists(group string) error
	AddUserToGroup(username string, group string) error
	RemoveUserFromGroup(username string, group string) error
}

type UserRepository interface {
	Create(username string, token string, group string) error
	GetToken(username string) (string, error)
	GetGroup(username string) (string, error)
	GetUsername(token string) (string, error)
	DeleteToken(username string) error
	DeleteGroup(username string) error
}

package repositories

import (
	chat "github.com/jayshrivastava/groupchat/proto"
)

type ChannelRepository interface {
	Open(key string) error
	Get(key string) (chan chat.StreamResponse, error)
	Close(key string) (error)
}

type GroupRepository interface {
	GetGroupMembers(group string, username string) ([]string, error)
	CreateIfNotExists(group string) error
	AddUserToGroup(username string, group string) error
	RemoveUserFromGroup(username string, group string) error
}

type UserRepository interface {
	Create(username string, token string, group string, password string) error
	GetToken(username string) (string, error)
	GetGroup(username string) (string, error)
	GetUsername(token string) (string, error)
	DeleteToken(username string) error
	DeleteGroup(username string) error
	CheckPassword(username string, candidatePassword string) (bool, error)
	DoesUserExist(username string) bool
	SetUserData(username string, token string, group string) error
}

package application

import (
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
)

type Context struct {
	ChannelRepository repositories.ChannelRepository
	GroupRepository   repositories.GroupRepository
	UserRepository    repositories.UserRepository
}

func CreateContext(
	channelRepository repositories.ChannelRepository,
	groupRepository repositories.GroupRepository,
	userRepository repositories.UserRepository,
) *Context {
	context := new(Context)
	context.ChannelRepository = channelRepository
	context.GroupRepository = groupRepository
	context.UserRepository = userRepository
	return context
}

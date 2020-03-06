package context

import (
	authentication "github.com/jayshrivastava/groupchat/server/authentication"
	repositories "github.com/jayshrivastava/groupchat/server/repositories"
)

type Context struct {
	ChannelRepository repositories.ChannelRepository
	GroupRepository   repositories.GroupRepository
	UserRepository    repositories.UserRepository
	Authenticator     authentication.Authenticator
}

func CreateContext(
	channelRepository repositories.ChannelRepository,
	groupRepository repositories.GroupRepository,
	userRepository repositories.UserRepository,
	authenticator authentication.Authenticator,
) *Context {
	context := new(Context)
	context.ChannelRepository = channelRepository
	context.GroupRepository = groupRepository
	context.UserRepository = userRepository
	context.Authenticator = authenticator
	return context
}

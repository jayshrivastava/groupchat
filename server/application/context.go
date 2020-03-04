package application

import(
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
)

type Context struct {
	ChannelRepository repositories.ChannelRepository
	GroupRepository repositories.GroupRepository
}

func CreateContext(
	channelRepository repositories.ChannelRepository,
	groupRepository repositories.GroupRepository,
) *Context {
	context := new(Context)
	context.ChannelRepository = channelRepository
	context.GroupRepository = groupRepository
	return context
}

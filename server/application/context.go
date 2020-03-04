package application

import(
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
)

type Context struct {
	ChannelRepository repositories.ChannelRepository
}

func CreateContext(
	channelRepository repositories.ChannelRepository,
) *Context {
	context := new(Context)
	context.ChannelRepository = channelRepository
	return context
}

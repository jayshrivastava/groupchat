package application

import(
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
)

func CreateApplicationContext() *Context{
	return CreateContext(
		repositories.CreateApplicationChannelRepository(), 
		repositories.CreateApplicationGroupRepository(),
	)
}


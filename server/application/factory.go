package application

import (
	authentication "github.com/jayshrivastava/groupchat/server/application/authentication"
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
	. "github.com/jayshrivastava/groupchat/server/context"
)

func CreateApplicationContext() *Context {
	userRepository := repositories.CreateApplicationUserRepository()
	return CreateContext(
		repositories.CreateApplicationChannelRepository(),
		repositories.CreateApplicationGroupRepository(),
		userRepository,
		authentication.CreateApplicationAuthenticator(userRepository),
	)
}

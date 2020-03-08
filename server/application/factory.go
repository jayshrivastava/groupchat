package application

import (
	authentication "github.com/jayshrivastava/groupchat/server/application/authentication"
	repositories "github.com/jayshrivastava/groupchat/server/application/repositories"
)

func CreateServer(serverPassword string, port string) *Server {
	userRepository := repositories.CreateApplicationUserRepository()
	server := Server{
		ServerPassword:    serverPassword,
		Port:              port,
		ChannelRepository: repositories.CreateApplicationChannelRepository(),
		GroupRepository:   repositories.CreateApplicationGroupRepository(),
		UserRepository:    userRepository,
		Authenticator:     authentication.CreateApplicationAuthenticator(userRepository),
	}

	return &server
}

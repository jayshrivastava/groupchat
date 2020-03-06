package authentication

import (
	repositories "github.com/jayshrivastava/groupchat/server/repositories"
)

func CreateApplicationAuthenticator(
	repository repositories.UserRepository,
) *ApplicationAuthenticator {
	return &ApplicationAuthenticator{userRepository: repository}
}

package authentication

import (
	"fmt"
	uuid "github.com/google/uuid"
	repositories "github.com/jayshrivastava/groupchat/server/repositories"
)

type ApplicationAuthenticator struct {
	userRepository repositories.UserRepository
}

func (authenticator *ApplicationAuthenticator) GenerateToken() string {
	return uuid.New().String()
}

func (authenticator *ApplicationAuthenticator) Authenticate(candidateToken string, candidateUsername string) bool {
	username, err := authenticator.userRepository.GetUsername(candidateToken)
	if err != nil {
		return false
	}
	if username != candidateUsername {
		return false
	}
	return true
}

func (authenticator *ApplicationAuthenticator) IsTokenValid(candidateToken string) bool {
	_, err := authenticator.userRepository.GetUsername(candidateToken)
	if err != nil {
		return false
	}
	return true
}

func (authenticator *ApplicationAuthenticator) DeleteToken(token string, username string) error {
	if !authenticator.Authenticate(token, username) {
		return fmt.Errorf("Invalid Username,Token pair %s,%s", username, token)
	}
	authenticator.userRepository.DeleteToken(token)
	return nil
}

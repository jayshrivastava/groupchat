package repositories

import (
	"fmt"
	"sync"
)

type ApplicationUserRepository struct {
	Users map[string]*user
	Tokens map[string]string
	RWMutex sync.RWMutex 
}

func (repository *ApplicationUserRepository) Create(username string, token string, group string) (error) {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; found {
		return fmt.Errorf("User %s username already exists", username)
	}

	repository.Users[username] = createUserData(token, group)
	repository.Tokens[token] = username

	return nil
}

func (repository *ApplicationUserRepository) GetToken(username string) (string, error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	if _, found := repository.Users[username]; !found {
		return "", fmt.Errorf("User %s username not found", username)
	}

	return repository.Users[username].Token, nil
}

func (repository *ApplicationUserRepository) GetGroup(username string) (string, error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	if _, found := repository.Users[username]; !found {
		return "", fmt.Errorf("User %s username not found", username)
	}

	return repository.Users[username].Group, nil
}

func (repository *ApplicationUserRepository) GetUsername(token string) (string, error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	if _, found := repository.Tokens[token]; !found {
		return "", fmt.Errorf("Invalid token. Could not look up username")
	}

	return repository.Tokens[token], nil
}

func (repository *ApplicationUserRepository) DeleteToken(username string) (error) {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; !found {
		return fmt.Errorf("User %s username not found", username)
	}

	delete(repository.Tokens, repository.Users[username].Token )
	repository.Users[username].Token = ""

	return nil
}

func (repository *ApplicationUserRepository) DeleteGroup(username string) (error) {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; !found {
		return fmt.Errorf("User %s username not found", username)
	}

	repository.Users[username].Group = ""

	return nil
}

/* Available in `repositories` package Only */
type user struct {
	Token string
	Group string
}

func createUserData(token string, group string) *user {
	return &user{token, group}
}

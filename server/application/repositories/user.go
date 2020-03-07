package repositories

import (
	"fmt"
	"sync"
)

type ApplicationUserRepository struct {
	Users   map[string]*user
	Tokens  map[string]string
	RWMutex sync.RWMutex
}

func (repository *ApplicationUserRepository) Create(username string, token string, group string, password string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; found {
		return fmt.Errorf("User %s username already exists", username)
	}

	repository.Users[username] = createUserData(token, group, password)
	repository.Tokens[token] = username

	return nil
}

func (repository *ApplicationUserRepository) SetUserData(username string, token string, group string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; !found {
		return fmt.Errorf("User %s username not found", username)
	}

	repository.Users[username].Token = token
	repository.Users[username].Group = group
	repository.Tokens[token] = username

	return nil
}

func (repository *ApplicationUserRepository) DoesUserExist(username string) bool {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	_, found := repository.Users[username]
	return found
}

func (repository *ApplicationUserRepository) CheckPassword(username string, candidatePassword string) (bool, error) {
	if !repository.DoesUserExist(username) {
		return false, fmt.Errorf("User %s does not exist", username)
	}
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	return repository.Users[username].Password == candidatePassword, nil
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

func (repository *ApplicationUserRepository) DeleteToken(username string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Users[username]; !found {
		return fmt.Errorf("User %s username not found", username)
	}

	delete(repository.Tokens, repository.Users[username].Token)
	repository.Users[username].Token = ""

	return nil
}

func (repository *ApplicationUserRepository) DeleteGroup(username string) error {
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
	Token    string
	Group    string
	Password string
}

func createUserData(token string, group string, password string) *user {
	return &user{token, group, password}
}

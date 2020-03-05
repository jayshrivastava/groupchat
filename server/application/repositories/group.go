package repositories

import (
	"fmt"
	"sync"
)

type ApplicationGroupRepository struct {
	Groups  map[string](map[string]bool)
	RWMutex sync.RWMutex
}

func (repository *ApplicationGroupRepository) GetGroupMembers(group string, username string) ([]string, error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	if _, found := repository.Groups[group]; !found {
		return make([]string, 0), fmt.Errorf("Group %s not found", group)
	}

	if _, found := repository.Groups[group][username]; !found {
		return make([]string, 0), fmt.Errorf("Username %s does not exist in %s", username, group)
	}

	usernames := make([]string, len(repository.Groups[group])-1)
	i := 0
	for memberUsername, _ := range repository.Groups[group] {
		if memberUsername != username {
			usernames[i] = memberUsername
			i += 1
		}
	}

	return usernames, nil
}

func (repository *ApplicationGroupRepository) CreateIfNotExists(group string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Groups[group]; !found {
		repository.Groups[group] = map[string]bool{}
	}

	return nil
}

func (repository *ApplicationGroupRepository) AddUserToGroup(username string, group string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Groups[group]; !found {
		return fmt.Errorf("Group %s not found", group)
	}

	if _, found := repository.Groups[group][username]; found {
		return fmt.Errorf("User %s already exists in group %s", username, group)
	}

	repository.Groups[group][username] = true

	return nil
}

func (repository *ApplicationGroupRepository) RemoveUserFromGroup(username string, group string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Groups[group]; !found {
		return fmt.Errorf("Group %s not found", group)
	}

	if _, found := repository.Groups[group][username]; !found {
		return fmt.Errorf("User %s does not exist in group %s", username, group)
	}

	delete(repository.Groups[group], username)

	return nil
}

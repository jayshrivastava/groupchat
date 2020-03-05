package repositories

import (
	"fmt"
	chat "github.com/jayshrivastava/groupchat/proto"
	"sync"
)

type ApplicationChannelRepository struct {
	Channels map[string]chan chat.StreamResponse
	RWMutex  sync.RWMutex
}

func (repository *ApplicationChannelRepository) Create(key string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

	if _, found := repository.Channels[key]; found {
		return fmt.Errorf("Duplicate key %s already exists", key)
	}

	repository.Channels[key] = make(chan chat.StreamResponse, 5)
	return nil
}

func (repository *ApplicationChannelRepository) Get(key string) (chan chat.StreamResponse, error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()
	if _, found := repository.Channels[key]; !found {
		return nil, fmt.Errorf("Could not find channel for key %s", key)
	}
	return repository.Channels[key], nil
}

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

func (repository *ApplicationChannelRepository) Open(key string) error {
	repository.RWMutex.Lock()
	defer repository.RWMutex.Unlock()

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

func (repository *ApplicationChannelRepository) Close(key string) (error) {
	repository.RWMutex.RLock()
	defer repository.RWMutex.RUnlock()

	if _, found := repository.Channels[key]; !found {
		return fmt.Errorf("Could not find channel for key %s", key)
	}
	close(repository.Channels[key])
	return nil
}


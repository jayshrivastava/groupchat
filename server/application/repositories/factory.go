package repositories

import (
	"sync"

	chat "github.com/jayshrivastava/groupchat/proto"
)

func CreateApplicationChannelRepository() *ApplicationChannelRepository {
	repo := new(ApplicationChannelRepository)
	repo.Channels = map[string]chan chat.StreamResponse{}
	repo.RWMutex = sync.RWMutex{}
	return repo
}

func CreateApplicationGroupRepository() *ApplicationGroupRepository {
	repo := new(ApplicationGroupRepository)
	repo.Groups = map[string](map[string]bool){}
	repo.RWMutex = sync.RWMutex{}
	return repo
}

func CreateApplicationUserRepository() *ApplicationUserRepository {
	repo := new(ApplicationUserRepository)
	repo.Users = map[string]*user{}
	repo.Tokens = map[string]string{}
	return repo
}

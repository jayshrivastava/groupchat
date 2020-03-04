package repositories

import(
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

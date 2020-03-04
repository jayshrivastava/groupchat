package repositories

import(
	"sync"
	
	chat "github.com/jayshrivastava/groupchat/proto"
)

func CreateApplicationRepository() *ApplicationChannelRepository {
	repo := new(ApplicationChannelRepository)
	repo.Channels = map[string]chan chat.StreamResponse{}
	repo.RWMutex = sync.RWMutex{}
	return repo
}
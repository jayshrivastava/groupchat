package repositories

import(
	chat "github.com/jayshrivastava/groupchat/proto"
)

type ChannelRepository interface {
	Create(key string) error
	Get(key string) (chan chat.StreamResponse, error)
}

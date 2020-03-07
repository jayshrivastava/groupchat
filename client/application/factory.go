package application

import (
	chat "github.com/jayshrivastava/groupchat/proto"
)

func CreateClient(
	rpc chat.ChatClient,
	username       string,
	userPassword   string,
	serverPassword string,
	group          string,
) *Client {
	server := Client{
		RPC: rpc,
		Username: username,
		UserPassword: userPassword,
		ServerPassword: serverPassword,
		Group: group,        
		Token: "",  
	}

	return &server
}
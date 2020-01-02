package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"flag"
	"syscall"
	"time"
	chat "./chat"
)

type ClientMeta struct {
	Username *string
	Password *string
	Host *string
	Group * string
	Token *string
}

func ClientError(e error) {
	fmt.Printf("%s\n", e.Error())
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func GetClientMeta() (ClientMeta) {
	cm := ClientMeta{
		Username: flag.String("u", "", "Username"),
		Password: flag.String("p", "", "Server Password"),
		Host: flag.String("h", "", "Host"),
		Group: flag.String("g", "", "Username"),
	}
	flag.Parse();

	if (*cm.Username == "" || *cm.Password == "" || *cm.Host == "" || *cm.Group == "") {
		ClientError(fmt.Errorf("Missing Flags"))
	}

	return cm
}

func Login(client chat.ChatClient, cm ClientMeta) {

	req := chat.LoginRequest{
		Username: *cm.Username,
		Password: *cm.Password,
		Group: *cm.Group,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.Login(ctx, &req)

	if err != nil {
		ClientError(fmt.Errorf("Login Failed: %s", err))
	}

}

func main() {

	// parse flags
	clientMeta := GetClientMeta()
	fmt.Printf("USERNAME %s\n", *clientMeta.Username)
	
	// register server
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := chat.NewChatClient(conn)

	// client login
	Login(client, clientMeta)

	// message sending thread

	// message receiving thread

	// register signal handler for logout
}
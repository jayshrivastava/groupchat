package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"flag"
	"syscall"
	"time"
	"os/signal"
	"sync"
	chat "./chat"
)

type ClientMeta struct {
	Username *string
	Password *string
	Host *string
	Group * string
	Token string
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

func Login(client chat.ChatClient, cm *ClientMeta) {

	req := chat.LoginRequest{
		Username: *cm.Username,
		Password: *cm.Password,
		Group: *cm.Group,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.Login(ctx, &req)

	if err != nil {
		ClientError(fmt.Errorf("Login Failed: %s", err))
	}

	cm.Token = res.Token

}

func Logout(client chat.ChatClient, cm *ClientMeta) {

	req := chat.LogoutRequest{
		Username: *cm.Username,
	}

	meta := metadata.New(map[string]string{"token": cm.Token})

	ctx := metadata.NewOutgoingContext(context.Background(), meta)
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()
	_, err := client.Logout(ctx, &req)

	if err != nil {
		ClientError(fmt.Errorf("Logout Failed: %s", err))
	}

}

func LogoutHandler(client chat.ChatClient, wg *sync.WaitGroup, cm *ClientMeta) {
	defer wg.Done()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	_ = <-sigs
	
	Logout(client, cm)
}

func main() {

	// parse flags
	clientMeta := GetClientMeta()
	
	// register server
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := chat.NewChatClient(conn)

	// client login
	Login(client, &clientMeta)

	// create waitgroup and dispatch threads
	wg := sync.WaitGroup{}
	wg.Add(1)
	
	// register signal handler for logout
	go LogoutHandler(client, &wg, &clientMeta)

	// message sending thread

	// message receiving thread

	wg.Wait()
}
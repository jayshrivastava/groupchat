package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	. "github.com/jayshrivastava/groupchat/helpers"
	chat "github.com/jayshrivastava/groupchat/proto"
)

type ClientMeta struct {
	Username string
	Password string
	Host     string
	Group    string
	Token    string
}

func Login(client chat.ChatClient, cm *ClientMeta) {

	req := chat.LoginRequest{
		Username: cm.Username,
		Password: cm.Password,
		Group:    cm.Group,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.Login(ctx, &req)

	if err != nil {
		Error(fmt.Errorf("Login Failed: %s", err))
	}

	cm.Token = res.Token

}

func Logout(client chat.ChatClient, cm *ClientMeta) {

	req := chat.LogoutRequest{
		Username: cm.Username,
	}

	meta := metadata.New(map[string]string{"token": cm.Token})

	ctx := metadata.NewOutgoingContext(context.Background(), meta)
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()
	_, err := client.Logout(ctx, &req)

	if err != nil {
		Error(fmt.Errorf("Logout Failed: %s", err))
	}

}

func LogoutHandler(client chat.ChatClient, wg *sync.WaitGroup, cm *ClientMeta) {
	defer wg.Done()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	_ = <-sigs

	Logout(client, cm)
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func Stream(client chat.ChatClient, wg *sync.WaitGroup, cm *ClientMeta) error {

	meta := metadata.New(map[string]string{"token": cm.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), meta)

	stream, err := client.Stream(ctx)

	if err != nil {
		Error(fmt.Errorf("Could not connect to stream: %s", err))
	}
	defer stream.CloseSend()

	go ClientSender(stream, cm)
	return ClientReceiver(stream, cm)
}

func ClientSender(stream chat.Chat_StreamClient, cm *ClientMeta) {

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		req := chat.StreamRequest{
			Username: cm.Username,
			Group:    cm.Group,
			Message:  text,
		}
		stream.Send(&req)
	}
}

func ClientReceiver(stream chat.Chat_StreamClient, cm *ClientMeta) error {

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			continue
		}

		switch evt := res.Event.(type) {
		case *chat.StreamResponse_ClientMessage:

			fmt.Printf("[%s] (%s) %s\n", TimestampToString(res.Timestamp), evt.ClientMessage.Username, evt.ClientMessage.Message)
		case *chat.StreamResponse_ClientLogin:
			fmt.Printf("[%s] (%s joined %s)\n", TimestampToString(res.Timestamp), evt.ClientLogin.Username, evt.ClientLogin.Group)
		case *chat.StreamResponse_ClientExisting:
			// timestamp exists but we do not need it
			fmt.Printf("(member %s of %s)\n", evt.ClientExisting.Username, evt.ClientExisting.Group)
		case *chat.StreamResponse_ClientLogout:
			fmt.Printf("[%s] (%s left %s)\n", TimestampToString(res.Timestamp), evt.ClientLogout.Username, evt.ClientLogout.Group)
		default:

		}
	}
}

func ClientMain(clientMeta ClientMeta) {

	// register server
	conn, err := grpc.Dial(clientMeta.Host, grpc.WithInsecure())
	if err != nil {
		Error(fmt.Errorf("fail to dial: %v", err))
	}
	defer conn.Close()
	client := chat.NewChatClient(conn)

	// client login
	Login(client, &clientMeta)

	// create waitgroup and dispatch threads
	wg := sync.WaitGroup{}

	// register signal handler for logout
	wg.Add(1)
	go LogoutHandler(client, &wg, &clientMeta)

	// streaming thread
	wg.Add(1)
	go Stream(client, &wg, &clientMeta)

	wg.Wait()
}

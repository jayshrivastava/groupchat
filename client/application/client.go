package application

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

	"google.golang.org/grpc/metadata"

	. "github.com/jayshrivastava/groupchat/helpers"
	chat "github.com/jayshrivastava/groupchat/proto"
)

type Client struct {
	RPC            chat.ChatClient
	Username       string
	UserPassword   string
	ServerPassword string
	Group          string
	Token          string
}

func (client *Client) Connect() {

	req := chat.LoginRequest{
		Username:       client.Username,
		UserPassword:   client.UserPassword,
		ServerPassword: client.ServerPassword,
		Group:          client.Group,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.RPC.Login(ctx, &req)

	if err != nil {
		Error(fmt.Errorf("Login Failed: %s", err))
	}

	client.Token = res.Token
}

func (client *Client) Logout() {

	req := chat.LogoutRequest{
		Username: client.Username,
	}

	meta := metadata.New(map[string]string{"token": client.Token})

	ctx := metadata.NewOutgoingContext(context.Background(), meta)
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()
	_, err := client.RPC.Logout(ctx, &req)

	if err != nil {
		Error(fmt.Errorf("Logout Failed: %s", err))
	}

}

func (client *Client) LogoutHandler(wg *sync.WaitGroup) {
	defer wg.Done()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	_ = <-sigs

	client.Logout()
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func (client *Client) Stream(wg *sync.WaitGroup) error {

	meta := metadata.New(map[string]string{"token": client.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), meta)

	stream, err := client.RPC.Stream(ctx)

	if err != nil {
		Error(fmt.Errorf("Could not connect to stream: %s", err))
	}
	defer stream.CloseSend()

	go client.sender(stream)
	return client.reciever(stream)
}

func (client *Client) sender(stream chat.Chat_StreamClient) {

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		req := chat.StreamRequest{
			Username: client.Username,
			Group:    client.Group,
			Message:  text,
		}
		stream.Send(&req)
	}
}

func (client *Client) reciever(stream chat.Chat_StreamClient) error {

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

func (client *Client) Run() {

	client.Connect()

	// create waitgroup and dispatch threads
	wg := sync.WaitGroup{}

	// register signal handler for logout
	wg.Add(1)
	go client.LogoutHandler(&wg)

	// streaming thread
	wg.Add(1)
	go client.Stream(&wg)

	wg.Wait()
}

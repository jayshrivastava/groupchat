package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/jayshrivastava/groupchat/helpers"
	clientApp "github.com/jayshrivastava/groupchat/client/application"
	serverApp "github.com/jayshrivastava/groupchat/server/application"
	chat "github.com/jayshrivastava/groupchat/proto"
)

type flags struct {
	Username       *string
	UserPassword   *string
	ServerPassword *string
	URL            *string
	Group          *string
	Port           *string
	RunAsServer    *bool
}

func main() {

	flags := flags{
		Username:       flag.String("u", "", "Username"),
		ServerPassword: flag.String("spass", "", "Server Password"),
		UserPassword:   flag.String("pass", "", "User Password"),
		URL:            flag.String("url", "", "Host"),
		Port:           flag.String("port", "5000", "Port"),
		Group:          flag.String("g", "", "Username"),
		RunAsServer:    flag.Bool("s", false, "Run server if flag is present"),
	}

	flag.Parse()

	if !*flags.RunAsServer {
		if *flags.Username == "" || *flags.UserPassword == "" || *flags.ServerPassword == "" || *flags.Group == "" || *flags.URL == "" {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		conn, err := grpc.Dial(*flags.URL, grpc.WithInsecure())
		if err != nil {
			helpers.Error(fmt.Errorf("fail to dial: %v", err))
		}
		defer conn.Close()
		rpcClient := chat.NewChatClient(conn)
		client := clientApp.CreateClient(
			rpcClient,
			 *flags.Username,
			 *flags.UserPassword,
			 *flags.ServerPassword,
			 *flags.Group,
		)
		client.Run()
	} else {

		if *(flags.ServerPassword) == "" || *(flags.Port) == "" {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *(flags.Port)))
		if err != nil {
			helpers.Error(fmt.Errorf("failed to listen: %v", err))
		}

		server := grpc.NewServer()
		chat.RegisterChatServer(server, serverApp.CreateServer(*(flags.ServerPassword), *(flags.Port)))
		server.Serve(lis)
	}
}

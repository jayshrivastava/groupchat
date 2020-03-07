package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/jayshrivastava/groupchat/client"
	"github.com/jayshrivastava/groupchat/helpers"
	"github.com/jayshrivastava/groupchat/server/application"
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
		cm := client.ClientMeta{
			Username:       *flags.Username,
			UserPassword:   *flags.UserPassword,
			ServerPassword: *flags.ServerPassword,
			Host:           *flags.URL,
			Group:          *flags.Group,
		}

		if cm.Username == "" || cm.UserPassword == "" || cm.Host == "" || cm.Group == "" || cm.ServerPassword == "" {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		client.ClientMain(cm)
	} else {

		if *(flags.ServerPassword) == "" || *(flags.Port) == "" {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *(flags.Port)))
		if err != nil {
			helpers.Error(fmt.Errorf("failed to listen: %v", err))
		}

		server := grpc.NewServer()
		chat.RegisterChatServer(server, application.CreateServer(*(flags.ServerPassword), *(flags.Port)))
		server.Serve(lis)
	}
}

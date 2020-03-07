package main

import (
	"flag"
	"fmt"

	"github.com/jayshrivastava/groupchat/client"
	"github.com/jayshrivastava/groupchat/helpers"
	"github.com/jayshrivastava/groupchat/server"
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
		sp := server.ServerProps{
			ServerPassword: *(flags.ServerPassword),
			Port:           *(flags.Port),
		}

		if sp.ServerPassword == "" || sp.Port == "" {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		server.ServerMain(sp)
	}
}

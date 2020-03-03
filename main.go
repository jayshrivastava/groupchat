package main

import (
	"flag"
	"fmt"

	"github.com/jayshrivastava/groupchat/server"
	"github.com/jayshrivastava/groupchat/client"
	"github.com/jayshrivastava/groupchat/helpers"
)

type Flags struct {
	Username *string
	Password *string
	Host *string
	Port *string
	Group *string
	RunAsServer *bool
}

func main() {

	flags := Flags{
		Username: flag.String("u", "", "Username"),
		Password: flag.String("p", "", "Server Password"),
		Host: flag.String("h", "", "Host"),
		Port: flag.String("port", "5000", "Port"),
		Group: flag.String("g", "", "Username"),
		RunAsServer: flag.Bool("s", false, "Run server if flag is present"),
	}

	flag.Parse();

	if (!*flags.RunAsServer) {
		cm := client.ClientMeta {
			Username: *flags.Username,
			Password: *flags.Password,
			Host: *flags.Host,
			Group: *flags.Group,
		}

		if (cm.Username == "" || cm.Password == "" || cm.Host == "" || cm.Group == "") {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}
	
		client.ClientMain(cm)
	} else {
		sp := server.ServerProps {
			Password: *(flags.Password),
			Port: *(flags.Port),
		}

		if (sp.Password == "" || sp.Port == "") {
			helpers.Error(fmt.Errorf("Missing Flags"))
		}

		server.ServerMain(sp)
	}
}
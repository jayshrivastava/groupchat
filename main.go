package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	Username *string
	Password *string
	Host *string
	Group *string
	RunAsServer *bool
}

func main() {

	flags := Flags{
		Username: flag.String("u", "", "Username"),
		Password: flag.String("p", "", "Server Password"),
		Host: flag.String("h", "", "Host"),
		Group: flag.String("g", "", "Username"),
		RunAsServer: flag.Bool("s", false, "Run server if flag is present"),
	}

	flag.Parse();

	if (!*flags.RunAsServer) {
		cm := ClientMeta{
			Username: flags.Username,
			Password: flags.Password,
			Host: flags.Host,
			Group: flags.Group,
		}

		if (*cm.Username == "" || *cm.Password == "" || *cm.Host == "" || *cm.Group == "") {
			Error(fmt.Errorf("Missing Flags"))
		}
	
		ClientMain(cm)
	} else {
		sp := ServerProps {
			Password: flags.Password,
			Host: flags.Host,
		}


		if (*sp.Password == "" || *sp.Host == "") {
			Error(fmt.Errorf("Missing Flags"))
		}

		ServerMain(sp)
	}
}
package main

import (
	"flag"
	"fmt"
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
		Port: flag.String("port", "", "Port"),
		Group: flag.String("g", "", "Username"),
		RunAsServer: flag.Bool("s", false, "Run server if flag is present"),
	}

	flag.Parse();

	if (!*flags.RunAsServer) {
		cm := ClientMeta {
			Username: *flags.Username,
			Password: *flags.Password,
			Host: *flags.Host,
			Group: *flags.Group,
		}

		if (cm.Username == "" || cm.Password == "" || cm.Host == "" || cm.Group == "") {
			Error(fmt.Errorf("Missing Flags"))
		}
	
		ClientMain(cm)
	} else {
		sp := ServerProps {
			Password: *(flags.Password),
			Port: *(flags.Port),
		}


		if (sp.Password == "" || sp.Port == "") {
			Error(fmt.Errorf("Missing Flags"))
		}

		ServerMain(sp)
	}
}
package main

import (
	"flag"
	"fmt"
)

func main() {

	cm := ClientMeta{
		Username: flag.String("u", "", "Username"),
		Password: flag.String("p", "", "Server Password"),
		Host: flag.String("h", "", "Host"),
		Group: flag.String("g", "", "Username"),
	}

	// cp := ServerProps{
	// 	Password: flag.String("p", "", "Server Password"),
	// 	Host: flag.String("h", "", "Host"),
	// }


	flag.Parse();

	if (*cm.Username == "" || *cm.Password == "" || *cm.Host == "" || *cm.Group == "") {
		ClientError(fmt.Errorf("Missing Flags"))
	}

	ClientMain(cm)
}
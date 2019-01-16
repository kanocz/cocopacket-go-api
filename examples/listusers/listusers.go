package main

// this example adds ips specified in command line (or stdin) to all slaves exist on master

import (
	"flag"
	"fmt"
	"log"
	"os"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url    = flag.String("url", "", "URL of cocopacket master instance")
	user   = flag.String("user", "", "username for authorization")
	passwd = flag.String("password", "", "password for authorization")
)

func main() {
	flag.Parse()
	if "" == *url {
		fmt.Println("Usage: ", os.Args[0], "[flags]")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	users, err := api.ListUsers()

	if nil != err {
		log.Fatalln("Error reading users list from master:", err)
	}

	for login, isAdmin := range users {
		if isAdmin {
			println("(a)", login)
		} else {
			println("(u)", login)
		}
	}
}

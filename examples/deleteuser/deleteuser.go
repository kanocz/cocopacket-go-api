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
	args := flag.Args()

	if "" == *url || 1 != len(args) {
		fmt.Println("Usage: ", os.Args[0], "[flags] login")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	users, err := api.DeleteUser(args[0])

	if nil != err {
		log.Fatalln("Error while remove user from master:", err)
	}

	for login, isAdmin := range users {
		if isAdmin {
			println("(a)", login)
		} else {
			println("(u)", login)
		}
	}
}

package main

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
	admin  = flag.Bool("admin", false, "add user as admin")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if "" == *url || 2 != len(args) {
		fmt.Println("Usage: ", os.Args[0], "[flags] login password")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	users, err := api.AddUser(args[0], args[1], *admin)

	if nil != err {
		log.Fatalln("Error adding user to master:", err)
	}

	for login, isAdmin := range users {
		if isAdmin {
			println("(a)", login)
		} else {
			println("(u)", login)
		}
	}
}

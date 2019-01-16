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
)

func main() {
	flag.Parse()
	args := flag.Args()

	if "" == *url || 1 != len(args) {
		fmt.Println("Usage: ", os.Args[0], "[flags] slave")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	err := api.DeleteSlave(args[0])

	if nil != err {
		log.Fatalln("Error while remove user from master:", err)
	}

	fmt.Println("OK")
}

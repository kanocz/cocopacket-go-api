package main

// this example list slaves configured on master

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kanocz/cocopacket-go-api"
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

	slaves, err := api.GetSlaveList()
	if nil != err {
		log.Fatalln("Error on ip list get:", err)
	}

	for _, slave := range slaves {
		println(slave)
	}
}

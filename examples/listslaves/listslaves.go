package main

// this example list slaves configured on master

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
	show   = flag.String("show", "name", "type of list, one of name, nameIP, namePort, IP, port")
)

func main() {
	flag.Parse()
	if "" == *url {
		fmt.Println("Usage: ", os.Args[0], "[flags]")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	if "name" == *show {
		slaves, err := api.GetSlaveList()
		if nil != err {
			log.Fatalln("Error on ip list get:", err)
		}

		for _, slave := range slaves {
			println(slave)
		}
		return
	}

	var (
		slaves map[string]string
		err    error
	)

	if "nameIP" == *show || "IP" == *show {
		slaves, err = api.GetSlavesIPs()
	} else {
		slaves, err = api.GetSlavesAddrs()
	}

	if nil != err {
		log.Fatalln("Error on ip list get:", err)
	}

	for slave, addr := range slaves {
		if "nameIP" == *show || "namePort" == *show {
			fmt.Print(slave + ": ")
		}
		fmt.Println(addr)
	}

}

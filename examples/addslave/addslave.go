package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url      = flag.String("url", "", "URL of cocopacket master instance")
	user     = flag.String("user", "", "username for authorization")
	passwd   = flag.String("password", "", "password for authorization")
	copyFrom = flag.String("copy", "", "name of slave to copy list ip to new created one")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if "" == *url || 3 != len(args) {
		fmt.Println("Usage: ", os.Args[0], "[flags] name ip port")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	port, err := strconv.ParseUint(args[2], 10, 16)
	if nil != err {
		log.Fatalln("Invalid port:", err)
	}
	if port < 1 || port > 65535 {
		log.Fatalln("Port can't be < 1 or > 65535 / ", port)
	}

	err = api.AddSlave(net.ParseIP(args[1]), uint16(port), args[0], *copyFrom)

	if nil != err {
		log.Fatalln("Error adding slave to master:", err)
	}

	fmt.Println("OK")
}

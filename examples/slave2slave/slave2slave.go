package main

// this example add ips found on "source slave" to "target slave"

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
	desc   = flag.String("description", "auto added ip", "description setted for new added ips")
	source = flag.String("source", "", "name of \"source slave\"")
	target = flag.String("target", "", "name of \"target slave\"")
)

func main() {
	var err error

	flag.Parse()
	if "" == *url || "" == *source || "" == *target {
		fmt.Println("Usage: ", os.Args[0], "[flags] -source XXX -target YYY")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	config, err := api.GetConfigInfo()
	if nil != err {
		log.Fatalln("Error loading config from master:", err)
	}

	slavesMap := map[string]bool{*target: true}

	var ips []string
	for ip, desc := range config.Ping.IPs {
		for _, s := range desc.Slaves {
			if *source == s {
				ips = append(ips, ip)
				break
			}
		}
	}

	if 0 == len(ips) {
		fmt.Println("No IPs on slave", *source)
		return
	}

	err = api.IPsSetSlaves(ips, slavesMap)
	if nil != err {
		log.Fatalln("Error set slaves for ips call:", err)
	}

	fmt.Println("OK")
}

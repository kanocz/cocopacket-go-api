package main

// this example adds ips specified in command line (or stdin) to all slaves exist on master

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url    = flag.String("url", "", "URL of cocopacket master instance")
	user   = flag.String("user", "", "username for authorization")
	passwd = flag.String("password", "", "password for authorization")
	desc   = flag.String("description", "auto added ip", "description setted for new added ips")
	group  = flag.String("group", "auto", "to which group add ips")
	fav    = flag.Bool("favorite", false, "Set new added ips as favorite")
	stdin  = flag.Bool("stdin", false, "read ip list from stdin instead of command line")
	slaves = flag.String("slaves", "", "comma separated list of slaves, all slaves if empty")
)

func main() {
	var err error

	flag.Parse()
	if "" == *url {
		fmt.Println("Usage: ", os.Args[0], "[flags] [ip[ ip[ ip...]]]")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	slaveList := strings.Split(*slaves, ",")
	if 0 == len(slaveList) || ((1 == len(slaveList)) && ("" == slaveList[0])) {
		slaveList, err = api.GetSlaveList()
		if nil != err {
			log.Fatalln("Error on slave list get:", err)
		}
	}

	var ips []string
	if *stdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ips = append(ips, strings.TrimSpace(scanner.Text()))
		}
	} else {
		ips = flag.Args()
	}

	// get current IP list to avoid adding just existing ips
	existIPs, err := api.GetConfigInfo()
	if nil != err {
		log.Fatalln("Error reading current config from master:", err)
	}

	ipsToAdd := make([]string, 0, len(ips))
	for _, ip := range ips {
		if _, exist := existIPs.Ping.IPs[ip]; !exist {
			ipsToAdd = append(ipsToAdd, ip)
		}
	}

	if 0 == len(ipsToAdd) {
		fmt.Println("No new IPs to add")
		return
	}

	err = api.AddIPs(ipsToAdd, slaveList, *desc, []string{*group + "->"}, *fav)
	if nil != err {
		log.Fatalln("Error add ips call:", err)
	}

	fmt.Println("OK")
}

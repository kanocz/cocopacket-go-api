package main

// this example adds ips specified in command line (or stdin) to all slaves exist on master

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kanocz/cocopacket-go-api"
)

var (
	url    = flag.String("url", "", "URL of cocopacket master instance")
	user   = flag.String("user", "", "username for authorization")
	passwd = flag.String("password", "", "password for authorization")
	stdin  = flag.Bool("stdin", false, "read ip list from stdin instead of command line")
)

func main() {
	flag.Parse()
	if "" == *url {
		fmt.Println("Usage: ", os.Args[0], "[flags] [ip[ ip[ ip...]]]")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

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

	ipsToRemove := make([]string, 0, len(ips))
	for _, ip := range ips {
		if _, exist := existIPs.Ping.IPs[ip]; exist {
			ipsToRemove = append(ipsToRemove, ip)
		}
	}

	if 0 == len(ipsToRemove) {
		fmt.Println("No existing IPs to remove")
		return
	}

	err = api.DeleteIPs(ipsToRemove)
	if nil != err {
		log.Fatalln("Error delete ips call:", err)
	}

	fmt.Println("OK")
}

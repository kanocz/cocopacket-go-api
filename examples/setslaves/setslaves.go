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
	stdin  = flag.Bool("stdin", false, "read ip list from stdin instead of command line")
	slaves = flag.String("slaves", "", "list of slaves in format +slave1,+slave2,-slave3 to add slave1, slave2 and remove slave3")
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
		log.Fatalln("No slaves givven")
	}

	slavesMap := map[string]bool{}
	for _, slave := range slaveList {
		if len(slave) < 2 {
			log.Fatalln("Slave", slave, "not in format +slave / -slave")
		}
		switch slave[0] {
		case '-':
			slavesMap[slave[1:]] = false
		case '+':
			slavesMap[slave[1:]] = true
		default:
			log.Fatalln("Slave", slave, "not in format +slave / -slave")
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

	if 0 == len(ips) {
		fmt.Println("Empty list of IPs")
		return
	}

	err = api.IPsSetSlaves(ips, slavesMap)
	if nil != err {
		log.Fatalln("Error set slaves for ips call:", err)
	}

	fmt.Println("OK")
}

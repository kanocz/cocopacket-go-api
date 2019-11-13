package main

// this example add or removes ips of group from/to specific slaves list

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url       = flag.String("url", "", "URL of cocopacket master instance")
	user      = flag.String("user", "", "username for authorization")
	passwd    = flag.String("password", "", "password for authorization")
	desc      = flag.String("description", "auto added ip", "description setted for new added ips")
	slaves    = flag.String("slaves", "", "list of slaves in format +slave1,+slave2,-slave3 to add slave1, slave2 and remove slave3")
	group     = flag.String("group", "auto", "group to adjust ips in")
	subgroups = flag.Bool("recursive", false, "include subgroups")
)

func main() {
	var err error

	flag.Parse()
	if "" == *url {
		fmt.Println("Usage: ", os.Args[0], "[flags]")
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

	err = api.GroupSetSlaves(*group, slavesMap, *subgroups)
	if nil != err {
		log.Fatalln("Error set slaves for ips in groups call:", err)
	}

	fmt.Println("OK")
}

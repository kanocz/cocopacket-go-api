package main

// this example adds ips specified in file (or stdin) to all slaves exist on master
// csv format is name_for_ip,category,subCategory,ip

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url      = flag.String("url", "", "URL of cocopacket master instance")
	user     = flag.String("user", "", "username for authorization")
	passwd   = flag.String("password", "", "password for authorization")
	filename = flag.String("filename", "stdin", "filename to read ip list csv from")
	slaves   = flag.String("slaves", "", "comma separated list of slaves, all slaves if empty")
)

func main() {
	var err error

	flag.Parse()
	if "" == *url {
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

	var scanner *bufio.Scanner

	if "stdin" == *filename {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatalf("failed opening file %s: %s", *filename, err)
		}
		scanner = bufio.NewScanner(file)
		defer file.Close()
	}

	ips := map[string]api.TestDesc{}

	for scanner.Scan() {
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ",")
		if nil == net.ParseIP(parts[3]) {
			log.Println("Invalid ip", parts[3])
			continue
		}
		if "" == parts[0] {
			parts[0] = parts[3]
		}

		ips[parts[3]] = api.TestDesc{
			Slaves:      slaveList,
			Description: parts[0],
			Groups:      []string{parts[1] + "->" + parts[2] + "->"},
		}
	}

	err = api.AddIPsRaw(ips)

	if nil != err {
		log.Fatalln("Error add ips call:", err)
	}

	fmt.Println("OK")
}

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
	url         = flag.String("url", "", "URL of cocopacket master instance")
	user        = flag.String("user", "", "username for authorization")
	passwd      = flag.String("password", "", "password for authorization")
	filename    = flag.String("filename", "stdin", "filename to read ip list csv from")
	slaves      = flag.String("slaves", "", "comma separated list of slaves, all slaves if empty")
	removeOther = flag.Bool("remove", false, "removes all IPs not listed in csv file (danger!)")
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
	oldList := map[string]api.TestDesc{}
	csvIPs := []string{}

	if *removeOther {
		config, _ := api.GetConfigInfo() // we can ignore error here - in worst case we'll delete no remaining IPs
		oldList = config.Ping.IPs
	}

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

		if *removeOther {
			csvIPs = append(csvIPs, parts[3])
		}
	}

	err = api.AddIPsRaw(ips)

	if nil != err {
		log.Fatalln("Error add ips call:", err)
	}

	if *removeOther {
		ipsToDelete := []string{}

		for ip := range oldList {
			if _, ok := ips[ip]; ok {
				continue
			}
			ipsToDelete = append(ipsToDelete, ip)
		}

		if len(ipsToDelete) > 0 {
			// remove IPs unlisted in csv
			err = api.DeleteIPs(ipsToDelete)
			if nil != err {
				log.Fatalln("Error delete ips call:", err)
			}

			// also remove
			allSlaves, err := api.GetSlaveList()
			if nil != err {
				log.Fatalln("Error receiving slaves call:", err)
			}

			newSlaves := map[string]bool{}
			for _, slave := range slaveList {
				newSlaves[slave] = true
			}
			counter := 0
			for _, slave := range allSlaves {
				if !newSlaves[slave] {
					newSlaves[slave] = false
					counter++
				}
			}
			if counter > 0 { // we have some slaves in system not listed in arguments, removeing from ips
				api.IPsSetSlaves(csvIPs, newSlaves)
			}
		}

	}

	fmt.Println("OK")
}

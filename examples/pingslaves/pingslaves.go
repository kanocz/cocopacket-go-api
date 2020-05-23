package main

/* this example adds slaves as ips to ping each other (using same uplink)
   to use it all slaves have to have name in format LOCATION-UPLINK like NYC-LEVEL3 (names like C-* are ignored)
   P.S.: this example show how to use raw ip adding */

import (
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
	group  = flag.String("group", "TRANSIT_PATHS", "to which group add ips")
	fav    = flag.Bool("favorite", false, "Set new added ips as favorite")
	dryrun = flag.Bool("dryrun", false, "just to show list what to add where, not actualy send to master")
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

	slaves, err := api.GetSlavesSources()
	if nil != err {
		log.Fatalln("Error on slave list get:", err)
	}

	uplinks := map[string][]string{}
	for slave := range slaves {
		parts := strings.Split(slave, "-")
		if "C" == parts[0] {
			continue
		}
		uplink := parts[len(parts)-1]
		uplinks[uplink] = append(uplinks[uplink], strings.Join(parts[:len(parts)-1], "-"))
	}

	ips := map[string]api.TestDesc{}

	for uplink, hosts := range uplinks {
		if len(hosts) > 1 {
			if *dryrun {
				fmt.Printf("%s: %+v\n", uplink, hosts)
			}
			for _, host := range hosts {
				// yes, it's better ways like append(s[:index], s[index+1:]...) exists, but this is just example :)
				xslaves := []string{}
				for _, s := range hosts {
					if host == s {
						continue
					}
					xslaves = append(xslaves, s+"-"+uplink)
				}

				if *dryrun {
					fmt.Printf("  %s (%s): %+v\n", host, slaves[host+"-"+uplink], xslaves)
				}

				// some slaves may not source have ip defined (unconnected, for example)
				if "" == slaves[host+"-"+uplink] {
					if *dryrun {
						fmt.Println(host+"-"+uplink, "doesn't have source IP defined!")
					}
					continue
				}

				ips[slaves[host+"-"+uplink]] = api.TestDesc{
					Slaves:      xslaves,
					Description: host,
					Groups:      []string{*group + "->" + uplink + "->"},
				}
			}
		}
	}

	if *dryrun {
		return
	}

	err = api.AddIPsRaw(ips)

	if nil != err {
		log.Fatalln("Error add ips call:", err)
		return
	}

	fmt.Println("OK")
}

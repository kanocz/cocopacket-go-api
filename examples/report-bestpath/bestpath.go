package main

/*
  this is very specific tool for one of bigger cocopacket clients, but maybe useful also as example
  all slaves in that system has names like NYC-DEFAULT, NYC-COGENT, NYC-LEVEL3 - so name of location and uplink
  purpose of this tool is detect/report situations when other uplink has better packet loss than DEFAULT
  __REPORT group name have to be used to get info for all selected for report ips
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url       = flag.String("url", "", "URL of cocopacket master instance")
	user      = flag.String("user", "", "username for authorization")
	passwd    = flag.String("password", "", "password for authorization")
	threshold = flag.Int("threshold", 3, "count of lost packet ignore (for example loss of 1-3 packets per hour may be ignored)")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if "" == *url || (1 != len(args) && 0 != len(args)) {
		fmt.Println("Usage: ", os.Args[0], "[flags] [groupname]")
		flag.Usage()
		return
	}

	groupname := "__report"
	if 1 == len(args) {
		groupname = args[0]
	}

	api.Init(*url, *user, *passwd)

	stats, err := api.GroupStats(groupname, true)

	if nil != err {
		log.Fatalln("Error reading data from master:", err)
	}

	// yes, we're testing only ping tests, http is not used by this client
	if 0 == len(stats.Ping) {
		fmt.Println("We got no ping data - please check report flags")
		return
	}

	// first we need re-group data to be able to compare them...
	// origin top keys is "ip@slave", we need to have ["ip"]["slave-location"]["slave-uplink"][timestamp]avgChunk

	data := map[string]*map[string]*map[string]map[int64]*api.AvgChunk{}

	for id, t := range stats.Ping {
		parts := strings.Split(id, "@")
		if 2 != len(parts) {
			continue // something goes wrong, skip is better solution :)
		}

		i := strings.LastIndex(parts[1], "-")
		if i < 0 {
			continue // we have no XXX-YYY scheme in slave name
		}

		ip := parts[0]
		slaveLocation := parts[1][:i]
		slaveUplink := parts[1][i+1:]

		dataIP, ok := data[ip]
		if !ok {
			dataIP = &map[string]*map[string]map[int64]*api.AvgChunk{}
			data[ip] = dataIP
		}

		dataIPLocation, ok := (*dataIP)[slaveLocation]
		if !ok {
			dataIPLocation = &map[string]map[int64]*api.AvgChunk{}
			(*dataIP)[slaveLocation] = dataIPLocation
		}

		(*dataIPLocation)[slaveUplink] = t
	}

	// we have all data, now justs compare DEFAULT to other for each ip and each timestamp

	for ip, locations := range data {
		ipHeaderPrinted := false // we pring IP header only if we have some interesting data

		for location, uplinks := range *locations {
			locationHeaderPrinted := false

			uDefault, ok := (*uplinks)["DEFAULT"]
			if !ok {
				continue // we want to have "-DEFAULT" uplink to compare
			}

			// first create timestamps slice to be able to sort it
			tss := make([]int64, 0, len(uDefault))
			for ts, s := range uDefault {
				if s.Loss < *threshold {
					// in case of no packet loss on default we don't need this timestamp in the future
					continue
				}
				tss = append(tss, ts)
			}
			sort.Slice(tss, func(i, j int) bool { return tss[i] < tss[j] })

			for _, ts := range tss {
				tsHeaderPrinted := false
				if 0 == uDefault[ts].Count || 0 == uDefault[ts].Loss {
					continue
				}
				defLoss := float64(uDefault[ts].Loss) / float64(uDefault[ts].Count) * 100

				for uname, udata := range *uplinks {
					if "DEFAULT" == uname {
						continue
					}

					chunk, ok := udata[ts]
					if !ok {
						continue
					}

					if 0 == chunk.Count {
						continue
					}

					uLoss := float64(chunk.Loss) / float64(chunk.Count) * 100

					if uLoss >= defLoss {
						continue
					}

					if !ipHeaderPrinted {
						ipHeaderPrinted = true
						fmt.Println(ip + ":   " + *url + "/#/detail?ip=" + ip + "&epoch=86400&probe=" + location + "-" + uname + "&graphType=basic\n")
					}

					if !locationHeaderPrinted {
						locationHeaderPrinted = true
						fmt.Println("   " + location)
					}

					if !tsHeaderPrinted {
						tsHeaderPrinted = true
						fmt.Print("     " + time.Unix(ts, 0).Format("2006-01-02 15:04") + " [ DEFAULT:" +
							strconv.FormatFloat(defLoss, 'f', 2, 64) + "% ")
					}

					fmt.Print(uname + ":" + strconv.FormatFloat(uLoss, 'f', 2, 64) + "% ")
				}

				if tsHeaderPrinted {
					fmt.Println("] ")
				}
			}

			// if locationHeaderPrinted {
			// 	fmt.Println("") // just make empty line
			// }

			if ipHeaderPrinted {
				fmt.Println("") // one more line :)
			}
		}

	}

}

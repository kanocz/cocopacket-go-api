package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url    = flag.String("url", "", "URL of cocopacket master instance")
	user   = flag.String("user", "", "username for authorization")
	passwd = flag.String("password", "", "password for authorization")
	report = flag.Bool("report", false, "limit only on tests selected for report using frontend")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if "" == *url || 1 != len(args) {
		fmt.Println("Usage: ", os.Args[0], "[flags] groupname")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	stats, err := api.GroupStats(args[0], *report)

	if nil != err {
		log.Fatalln("Error reading data from master:", err)
	}

	// yes, it's copy-paste here... but its just for not so good example :)
	if len(stats.HTTP) > 0 {
		fmt.Println("We have some HTTP data in this group! :)")
		for id, t := range stats.HTTP {
			fmt.Println(" * " + id)
			for ts, data := range t {
				timestamp := time.Unix(ts, 0)
				fmt.Print("    [" + timestamp.Format(time.RFC822) + "] ")
				lossStr := ""
				latStr := ""

				if data.Loss > 0 {
					lossStr = strconv.FormatFloat(float64(data.Loss)/float64(data.Count)*100, 'f', 2, 64) + "% loss"
				} else {
					lossStr = ""
				}

				if data.Latency > 0 {
					latStr = strconv.FormatFloat(float64(data.Latency)/float64(data.Count), 'f', 2, 64) + " ms"
				}

				fmt.Println(latStr, lossStr)
			}
		}
	}

	if len(stats.Ping) > 0 {
		fmt.Println("We have some PING data in this group! :)")
		for id, t := range stats.Ping {
			fmt.Println(" * " + id)
			for ts, data := range t {
				timestamp := time.Unix(ts, 0)
				fmt.Print("    [" + timestamp.Format(time.RFC822) + "] ")
				lossStr := ""
				latStr := ""

				if data.Loss > 0 {
					lossStr = strconv.FormatFloat(float64(data.Loss)/float64(data.Count)*100, 'f', 2, 64) + "% loss"
				} else {
					lossStr = ""
				}

				if data.Latency > 0 {
					latStr = strconv.FormatFloat(float64(data.Latency)/float64(data.Count), 'f', 2, 64) + " ms"
				}

				fmt.Println(latStr, lossStr)
			}
		}

	}
}

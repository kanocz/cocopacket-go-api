package main

// example of cocopacket to prometheus connector
// handling requests like http://127.0.0.1:8008/GROUP/SLAVE

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	api "github.com/kanocz/cocopacket-go-api"
)

var (
	url    = flag.String("url", "", "URL of cocopacket master instance")
	user   = flag.String("user", "", "username for authorization")
	passwd = flag.String("password", "", "password for authorization")
	report = flag.Bool("report", false, "limit only on tests selected for report using frontend")
	listen = flag.String("listen", "0.0.0.0:8008", "ip:port to listen for prometheus requests")
	extra  = flag.String("extra", ", department=\"cocopacket\"", "additional string metrics")
)

type prometheusHandler struct{}

func (prometheusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.RequestURI(), "/"), "/")
	if 2 != len(parts) {
		// invalid URL - we need /GROUP/SLAVE
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.Header().Add("Content-type", "text/plain")

	// ignoring second part, we don't need stats about http requests
	pings, _, err := api.GroupLastStats(parts[0], parts[1])
	if nil != err {
		w.Write([]byte("# error loading cocopacket stats: " + err.Error()))
		return
	}

	w.Write([]byte("# HELP cocopacket_latency in ms\n# TYPE cocopacket_latency gauge\n"))
	for ip, data := range pings {
		w.Write([]byte(fmt.Sprintf("cocopacket_latency{ip=\"%s\"%s} %.2f\n", ip, *extra, data.Latency)))
	}

	w.Write([]byte("# HELP cocopacket_packet_loss_percent\n# TYPE cocopacket_packet_loss_percent gauge\n"))
	for ip, data := range pings {
		if 0 == data.Loss {
			w.Write([]byte(fmt.Sprintf("cocopacket_packet_loss_percent{ip=\"%s\"%s} 0\n", ip, *extra)))
		} else {
			w.Write([]byte(fmt.Sprintf("cocopacket_packet_loss_percent{ip=\"%s\"%s} %.2f\n", ip, *extra, float64(data.Loss)/float64(data.Count)*100)))
		}
	}
}

func main() {
	flag.Parse()

	if "" == *url {
		fmt.Println("Please specify URL of cocopacket master")
		flag.Usage()
		return
	}

	api.Init(*url, *user, *passwd)

	srv := http.Server{
		Handler: prometheusHandler{},
		Addr:    fmt.Sprintf(*listen),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		// sig is a ^C, handle it
		fmt.Println("shutting down..")

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// start http shutdown
		srv.Shutdown(ctx)

		// verify, in worst case call cancel via defer
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("not all connections done")
			os.Exit(0)
		case <-ctx.Done():

		}
	}()

	log.Println(srv.ListenAndServe())
}

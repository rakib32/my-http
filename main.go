package main

import (
	"flag"
	"fmt"
	"my-http/httphelper"
	"my-http/worker"
	"net/http"
	neturl "net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	parallel   = flag.Int("parallel", 10, "the numbers of parallel request")
	dispatcher *worker.Dispatcher
)

type webScrapper struct {
	client *http.Client
}

func (ws *webScrapper) PullContent(url string) (string, error) {
	var md5String string

	md5String, err := httphelper.Get(url, ws.client)
	if err != nil {
		return "", err
	}

	return md5String, nil
}

func init() {
	flag.Parse()
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGQUIT)

	//Init processReqDispatcher
	dispatcher = worker.Init(200, *parallel)
	processRequests()

	<-stop

	dispatcher.Stop()
	fmt.Println("Shutdown Process ...")
}

func processRequests() {
	urls := flag.Args()

	for _, url := range urls {
		u, err := neturl.Parse(url)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		if u.Scheme == "" {
			u.Scheme = "http"
		}
		dispatcher.AddJob(worker.Job{
			Payload: u.String(),
			Handler: handleRequest(),
		})
	}
}

var handleRequest = func() worker.Handler {
	return func(url interface{}) (interface{}, error) {
		client := httphelper.CustomClient()
		scrapper := webScrapper{client: client}

		md5String, err := scrapper.PullContent(fmt.Sprintf("%v", url))
		if err != nil {
			return "", err
		}

		fmt.Printf("%v %v \n", url, md5String)
		return url, nil
	}
}

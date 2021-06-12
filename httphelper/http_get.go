package httphelper

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string, client *http.Client) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func CustomClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 20
	transport.MaxConnsPerHost = 20
	transport.MaxIdleConnsPerHost = 20

	client := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &client
}

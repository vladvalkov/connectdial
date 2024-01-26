package main

import (
	"github.com/elazarl/goproxy"
	"github.com/vladvalkov/connectdial"
	"log"
	"net/http"
)

func main() {
	createProxyServer()

	dialer, err := connectdial.New(connectdial.Config{
		ConnectionString: "127.0.0.1:8888",
	})
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{DialContext: dialer.DialContext}
	client := http.Client{Transport: transport}

	resp, err := client.Get("https://example.com")
	if err != nil {
		panic(err)
	}

	log.Println(resp.StatusCode)
}

func createProxyServer() {
	go func() {
		err := http.ListenAndServe(":8888", goproxy.NewProxyHttpServer())
		if err != nil {
			panic(err)
		}
	}()
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
)

var port string

// var originServerURL *url.URL
var originServerURL string

func main() {
	if len(os.Args) < 2 {
		port = ":8080"
	} else {
		port = os.Args[1]
	}
	fmt.Println("Port", port)

	originServerURL = "http://localhost:8080"

	//Establish a socket connection
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error on server start: %v\n", err)
		listener.Close()
		os.Exit(1)
	}

	//Handlers
	http.HandleFunc("/", proxyHandler)
	//Serve http requests that come through
	if err = http.Serve(listener, nil); err != nil {
		fmt.Printf("Error serving HTTP requests: %v\n", err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("File-path: %s\n", r.URL)
	fmt.Printf("File-name: %s\n", path.Base(r.URL.Path))

	if !checkMethod(w, r) {
		return
	}
	getHandler(w, r)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------- <Sending To Server> ----------")
	resp, err := http.Get(originServerURL + r.URL.String())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("---------- <Receiving From Server> ----------")

	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	
	fmt.Println("---------- <Sending To Client> ----------")
}

// Helper Methods
func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	fmt.Printf("Checked Method: %s\n", r.Method)
	switch r.Method {
	case "GET":
		return true
	default:
		http.Error(w, "Request Method Is Currently Not Supported", 501)
		return false
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

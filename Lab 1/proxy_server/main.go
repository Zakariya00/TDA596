package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

var port string
var forwardToServerURL string

func main() {
	//set up port & server forward address 
	if len(os.Args) < 3 {
		os.Exit(1)
	} else {
		port = os.Args[1]
		forwardToServerURL = os.Args[2]
	}
	
	fmt.Println("Port", port)
	fmt.Println("Forwarding address: ", forwardToServerURL)
	
	//Establish a socket connection
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error on server start: %v\n", err)
		listener.Close()
		os.Exit(1)
	}

	//Handlers
	http.HandleFunc("/", proxyHandler)

	//Listens for and serves http request by spawning a go process with handler
	if err = http.Serve(listener, nil); err != nil {
		fmt.Printf("Error serving HTTP requests: %v\n", err)
		os.Exit(1)
	}
}

// Main handler for proxy server, calls checkMethod and getHandler methods
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r) {
		return
	}
	getHandler(w, r)
}

// Serves get requests by forwarding them to the Main server and sends back response
func getHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(forwardToServerURL + r.URL.String())
	if err != nil {
		http.Error(w, "Failure in forwarding request: "+err.Error(), 400)
		return
	}

	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error in forwarding main server response: "+err.Error(), 500)
		return
	}
}

// Check for valid Http Requesting Method
func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	switch r.Method {
	case "GET":
		return true
	default:
		http.Error(w, "Request Method Is Currently Not Supported <"+r.Method+">", 501)
		return false
	}
}

// Copy source headers to destination header (use for sending back main server response)
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

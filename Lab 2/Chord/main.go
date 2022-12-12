package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var cNode = ChordNode{}

func main() {
	commandlineFlags()

	err := rpc.Register(&cNode)
	if err != nil {
		log.Fatal("error registering node", err)
	}
	go rpc.HandleHTTP()
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", serverHandler)

	l, e := net.Listen("tcp", ":"+cNode.LocalNode.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	for {
		cNode.handleKeyBoard()
	}
}

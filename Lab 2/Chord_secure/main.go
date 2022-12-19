package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var cNode = ChordNode{}

func main() {
	fmt.Printf("--------------------- <%s> ---------------------\n", getLocalAddress())
	commandlineFlags()

	err := rpc.Register(&cNode)
	if err != nil {
		log.Fatal("error registering node", err)
	}
	rpc.HandleHTTP()
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", serverHandler)

	l, e := net.Listen("tcp", ":"+cNode.LocalNode.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go func() {
		if err := http.Serve(l, nil); err != nil {
			fmt.Printf("Error serving HTTP requests: <%v>\n", err)
			os.Exit(1)
		}
	}()

	for {
		cNode.handleKeyBoard()
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var cNode = ChordNode{}

func main() {
	fmt.Printf("--------------------- <%s> ---------------------\n", getLocalAddress())
	commandlineFlags()

	err := rpc.Register(&cNode)
	if err != nil {
		log.Fatal("error registering node", err)
	}

	l, e := net.Listen("tcp", ":"+cNode.LocalNode.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	newP, e1 := strconv.Atoi(cNode.LocalNode.Port)
	l1, e := net.Listen("tcp", ":"+strconv.Itoa(newP+1))
	if e != nil || e1 != nil {
		log.Fatal("listen error:", e)
	}
	srv, err := httpsServer("client.crt")
	if err != nil {
		log.Fatal("Failed to set https Server:", err)
	}
	srv.Addr = ":" + strconv.Itoa(newP+1)

	go func() {
		rpc.HandleHTTP()
		if err := http.Serve(l, nil); err != nil {
			fmt.Printf("Error serving HTTP requests: <%v>\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		http.HandleFunc("/favicon.ico", faviconHandler)
		http.HandleFunc("/", serverHandler)
		if err := srv.ServeTLS(l1, "secure_chord.crt", "secure_chord.key"); err != nil {
			fmt.Printf("Error serving HTTP requests: <%v>\n", err)
			os.Exit(1)
		}
	}()

	for {
		cNode.handleKeyBoard()
	}
}

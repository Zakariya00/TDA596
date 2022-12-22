package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

/* Helper Functions */

// getLocalAddress finds local ip address
func getLocalAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// fileExists checks if file exists locally
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	} else {
		fmt.Printf("File does not exist\n")
		return false
	}
}

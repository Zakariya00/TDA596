package main

import (
	"flag"
	"fmt"
	"time"
)

/* Function for handling commandline flag args */

func commandlineFlags() {
	ip1 := flag.String("a", "", "The IP address that the Chord client will bind to")
	port1 := flag.String("p", "", "The port that the Chord client will bind to and listen on")
	ip2 := flag.String("ja", "", "The IP address of the machine running a Chord node")
	port2 := flag.String("jp", "", "The port that an existing Chord node is bound to and listening on")
	delay1 := flag.Int("ts", 30000, "The time in milliseconds between invocations of ‘stabilize")
	delay2 := flag.Int("tff", 10000, "The time in milliseconds between invocations of ‘fix fingers’")
	delay3 := flag.Int("tcp", 40000, "The time in milliseconds between invocations of ‘check predecessor")
	delay4 := flag.Int("s", 1, "The time in minutes between invocations of ‘backupHandler")
	nbrSuccesors := flag.Int("r", 3, "The number of successors maintained by the Chord client")
	idOverwrite := flag.String("i", "", "The identifier (ID) assigned to the Chord client which will"+
		" override the ID computed by the SHA1 sum of the client’s IP address and port number")
	debuggingOn := flag.Bool("d", false, "The switch for debugging print")

	flag.Parse()

	if *ip1 == "" {
		fmt.Println("Local Node IP hasn't been set\n<Setting to default (localhost IP)>")
		*ip1 = getLocalAddress()
	}
	if *port1 == "" {
		fmt.Println("Local Node Port hasn't been set\n<Setting to default (8080)>")
		*port1 = "8080"
	}

	debugging = *debuggingOn
	if *delay1 > 60000 || *delay1 < 1 {
		*delay1 = 30000
	}
	if *delay2 > 60000 || *delay2 < 1 {
		*delay2 = 10000
	}
	if *delay3 > 60000 || *delay3 < 1 {
		*delay3 = 40000
	}
	if *delay4 > 10080 || *delay4 < 1 {
		*delay4 = 1
	}
	if *nbrSuccesors > 32 || *nbrSuccesors < 1 {
		*nbrSuccesors = 3
	}
	stabilizationDelay = time.Duration(*delay1)
	fixFingersDelay = time.Duration(*delay2)
	predeccesorCheckDelay = time.Duration(*delay3)
	backupTimeDelay = time.Duration(*delay4)
	m = *nbrSuccesors

	ip := *ip1
	port := *port1
	address := ip + ":" + port
	id := hash(address).String()
	if (*idOverwrite != "") && (len(*idOverwrite) == 48) {
		id = *idOverwrite
	}
	cNode.LocalNode = &Node{id, ip, port, address}
	fmt.Printf("<LocalNode>: %+v\n", *cNode.LocalNode)

	if *ip2 != "" || *port2 != "" {
		hostIP := *ip2
		hostPort := *port2
		hostAddress := hostIP + ":" + hostPort
		hostID := hash(hostAddress).String()
		hostNode := Node{hostID, hostIP, hostPort, hostAddress}
		fmt.Printf("<HostNode>: %+v\n", hostNode)
		cNode.join(&hostNode)
	} else {
		cNode.create()
	}
}

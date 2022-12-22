package main

import (
	"fmt"
	"strings"
	"time"
)

/* User Interface and Input Functions */

func (chord *ChordNode) handleKeyBoard() {
	fmt.Println("\nType In The Desired Command:")
	var operation string
	fmt.Scan(&operation)

	switch strings.ToLower(operation) {
	case "lookup":
		fmt.Println("Input the file you want to search for")
		var file string
		fmt.Scan(&file)
		chord.lookUp(file)

	case "storefile":
		fmt.Println("Input the file you want to store")
		var file string
		fmt.Scan(&file)
		chord.storeFile(file)

	case "deletefile":
		fmt.Println("Enter the file you wish to delete")
		var file string
		fmt.Scan(&file)
		chord.deleteFile(file)

	case "ping":
		fmt.Println("Enter The Receiving Address For The Ping <IP:Port>: ")
		var address string
		fmt.Scan(&address)
		reply, err := chord.call(address, "ChordNode.Ping",
			RpcArgs{"0", "Ping!", nil, nil, nil})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(reply.Value)

	case "printstate":
		chord.printState()

	case "quit":
		chord.quit()

	case "debug":
		chord.debug()

	case "hash":
		fmt.Println("Enter the string you wish to hash:")
		var stringHash string
		fmt.Scan(&stringHash)
		fmt.Println(hash(stringHash))

	case "find":
		fmt.Println("Enter the key you wish to find the successor for:")
		var key string
		fmt.Scan(&key)
		has := chord.find(chord.LocalNode, key)
		fmt.Printf("found: %+v\n", *has)

	case "whoami":
		fmt.Printf("Node: %+v\n", *chord.LocalNode)

	case "sendtest":
		fmt.Println("Which file?:")
		var file string
		fmt.Scan(&file)
		if fileExists(file) {
			fmt.Println("Receiving address?")
			var address string
			fmt.Scan(&address)
			go postSender(address, file)
		}

	case "save":
		chord.backupFiles()

	case "getbackups":
		chord.movBackups()

	case "sysinfo":
		fmt.Println("Succession List Size:", m)
		fmt.Println("Stabilize delay time (ms):", stabilizationDelay*time.Millisecond)
		fmt.Println("Fix_fingers delay time (ms):", fixFingersDelay*time.Millisecond)
		fmt.Println("check_predecessor delay time (ms):", predeccesorCheckDelay*time.Millisecond)
		fmt.Println("backupHandler delay time (min):", backupTimeDelay*time.Minute)

	default:
		fmt.Println("\nCommand Is Not <Supported>/<Faulty Input>")
		fmt.Printf("Supported Commands are:\n- LookUp <file name>\n- StoreFile <file path>\n" +
			"- PrintState <>\n- Deletefile <file name>\n- Quit <>\n\nFor Debugging:\n- Ping <address>\n- Find <key>\n" +
			"- Hash <string>\n- Whoami <>\n- Sysinfo <>\n- sendtest <file name> <address>\n- Save <>\n- GetBackUps <>\n- Debug <>\n")
	}
}

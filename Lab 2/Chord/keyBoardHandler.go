package main

import (
	"fmt"
	"strings"
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
		chord.LookUp(file)

	case "storefile":
		fmt.Println("Input the file you want to store")
		var file string
		fmt.Scan(&file)
		chord.StoreFile(file)

	case "deletefile":
		fmt.Println("Enter the file you wish to delete")
		var file string
		fmt.Scan(&file)
		chord.DeleteFile(file)

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
		chord.PrintState()

	case "quit":
		chord.Quit()

	case "debug":
		chord.Debug()

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

	default:
		fmt.Println("\nCommand Is Not <Supported>/<Faulty Input>")
		fmt.Printf("Supported Commands are:\n- LookUp <file name>\n- StoreFile <file path>\n" +
			"- PrintState <>\n- Quit <>\n\n" +
			"For Debugging:\n- Ping <address>\n- Find <key>\n- Hash <string>\n- Whoami <>\n- Debug <>\n")
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/* Supported User Operations Functions */

var debugging bool

// LookUp look up file and return the node it belongs to
func (chord *ChordNode) LookUp(fileName string) *Node {
	fileKey := hash(fileName).String()
	hasIt := chord.find(chord.LocalNode, fileKey) //
	reply, _ := chord.call(hasIt.Address, "ChordNode.Get", RpcArgs{fileKey,
		"", nil, nil, nil})

	fmt.Printf("FileName <%s>\nFileKey <%s>\n", fileName, fileKey)
	fmt.Printf("File Should Be At Node: %+v\nFile Path: %s\n", *hasIt, reply.Value)
	return hasIt
}

// StoreFile send local file to the node it belongs to for storage
func (chord *ChordNode) StoreFile(filePath string) *Node {
	if !fileExists(filePath) {
		return nil
	}
	_, fileName := filepath.Split(filePath)
	fileKey := hash(fileName).String()
	sendTo := chord.LookUp(fileName)
	postSender(sendTo.Address, filePath)
	chord.call(sendTo.Address, "ChordNode.Put", RpcArgs{fileKey,
		fileName, nil, nil, nil})

	fmt.Printf("File Sent To Node: %+v\nNew File Path: %s\n", *sendTo, fileName)
	return sendTo
}

// DeleteFile delete the file from the node it should be at
func (chord *ChordNode) DeleteFile(filePath string) *Node {
	fileKey := hash(filePath).String()
	sendTo := chord.LookUp(filePath)
	chord.call(sendTo.Address, "ChordNode.Delete", RpcArgs{fileKey, "",
		nil, nil, nil})

	fmt.Printf("Deleted File From Node: %+v\n", *sendTo)
	return sendTo
}

// PrintState print out the current state of local, successors and Fingertable nodes
func (chord *ChordNode) PrintState() {
	hours, mins, secs := time.Now().Clock()
	fmt.Printf("--------------- <%d:%d:%d> ---------------\n", hours, mins, secs)
	if chord.Predecessor != nil {
		fmt.Printf("Predecessor: %+v\n", *chord.Predecessor)
	} else {
		fmt.Printf("Predecessor: {nil}\n")
	}
	fmt.Printf("LocalNode  : %+v\n", *chord.LocalNode)

	for i := 0; i < m; i++ {
		fmt.Printf("Succesor[%d]: %+v\n", i, *chord.Successor[i])
	}

	fmt.Println("Stored <key,value> pairs:")
	for key, value := range chord.Bucket {
		fmt.Printf("K: <" + key + "> -> " + "V: <" + value + ">\n")
	}

	for i := 0; i < keySize; i++ {
		fmt.Printf("FingerTable[%s]: %+v\n", chord.FingerTable[i].Start, *chord.FingerTable[i].Successor)
	}

}

// Quit print final state and send any keys to successor before exiting the program
func (chord *ChordNode) Quit() {
	fmt.Println("Exit Protocol Engaged: <Printing Final State>")
	chord.PrintState()
	chord.put_all()
	os.Exit(0)
}

func (chord *ChordNode) Debug() {
	if debugging {
		debugging = false
		fmt.Println("Debugging Turned Off")
		return
	}
	debugging = true
	fmt.Println("Debugging Turned On")
}

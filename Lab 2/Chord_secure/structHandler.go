package main

/* Chord Datastructures */

var m int

// Node representing a simple node
type Node struct {
	Id      string
	Ip      string
	Port    string
	Address string
}

// Entry representing an entry in the FingerTable
type Entry struct {
	Start     string
	Successor *Node
}

// ChordNode representing a chord node
type ChordNode struct {
	LocalNode   *Node
	FingerTable []*Entry
	Successor   []*Node
	Predecessor *Node

	Bucket  map[string]string
	Backups map[string]string
}

// RpcArgs for use as request and response between rpc caller and receiver
type RpcArgs struct {
	Key    string            // Key for passing/receiving key arg, ex key, Id
	Value  string            // Value for passing/receiving value/msg arg, ex Bucket[key]
	RNode  *Node             // RNode for passing/receiving node arg, ex LocalNode, predecessor or Successor[x]
	RNodes []*Node           // RNodes for passing/receiving nodes slice arg, ex Successor slice
	Keys   map[string]string // Keys for passing/receiving maps, ex Bucket
}

type handler struct{}

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

	Bucket map[string]string
}

// RpcArgs for use as request and response between rpc caller and receiver
type RpcArgs struct {
	Key    string
	Value  string
	RNode  *Node
	RNodes []*Node
	Keys   map[string]string
}

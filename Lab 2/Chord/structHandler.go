package main

/* Chord Datastructures */

var m int

type Node struct {
	Id      string
	Ip      string
	Port    string
	Address string
}

type Entry struct {
	Start     string
	Successor *Node
}

type ChordNode struct {
	LocalNode   *Node
	FingerTable []*Entry
	Successor   []*Node
	Predecessor *Node

	Bucket map[string]string
}

type RpcArgs struct {
	Key    string
	Value  string
	RNode  *Node
	RNodes []*Node
	Keys   map[string]string
}

package main

import (
	"fmt"
)

/* Ring Create/Join Functions*/

// create  a new chord ring
func (chord *ChordNode) create() {
	fmt.Println("<Create Chord Ring>")
	chord.Predecessor = nil
	chord.FingerTable = new([keySize + 2]*Entry)[1:(keySize + 1)]

	for i := 0; i < keySize; i++ {
		chord.FingerTable[i] = &Entry{}
		chord.FingerTable[i].Start = jump(chord.LocalNode.Address, i+1).String()
		chord.FingerTable[i].Successor = chord.LocalNode
	}

	chord.Successor = make([]*Node, m)
	for i := 0; i < m; i++ {
		chord.Successor[i] = chord.LocalNode
	}
	chord.Bucket = make(map[string]string)
	chord.Backups = make(map[string]string)

	/* Activate Background Processes*/
	chord.backGroundProcesses()
}

// join an existing chord ring
func (chord *ChordNode) join(hostNode *Node) {
	fmt.Println("<Join Chord Ring>")
	chord.Predecessor = nil
	chord.Successor = make([]*Node, m)

	// depending on the set stabilize & fix_fingers delay time,
	//might get back myself when joining ring after just leaving
	found := chord.find(hostNode, chord.LocalNode.Id)
	if found.Id == chord.LocalNode.Id {
		chord.Successor[0] = hostNode
	} else {
		chord.Successor[0] = found
	}

	chord.Backups = make(map[string]string)
	chord.Bucket = chord.get_all()
	chord.FingerTable = new([keySize + 2]*Entry)[1:(keySize + 1)]
	for i := 0; i < keySize; i++ {
		chord.FingerTable[i] = &Entry{}
		chord.FingerTable[i].Start = jump(chord.LocalNode.Address, i+1).String()
		chord.FingerTable[i].Successor = chord.Successor[0]
	}

	/* Activate Background Processes*/
	chord.backGroundProcesses()
}

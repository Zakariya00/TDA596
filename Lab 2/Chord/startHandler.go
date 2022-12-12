package main

import (
	"fmt"
)

/* Ring Create/Join Functions*/

// Functions for Creating or Joining rings
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

	/* Activate Background Processes*/
	chord.backGroundProcesses()
}

func (chord *ChordNode) join(hostNode *Node) {
	fmt.Println("<Join Chord Ring>")
	chord.Predecessor = nil
	chord.Successor = make([]*Node, m)
	chord.Successor[0] = chord.find(hostNode, chord.LocalNode.Id)
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

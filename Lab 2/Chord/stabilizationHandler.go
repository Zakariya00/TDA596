package main

import (
	"fmt"
)

/* Ring Stabilizing Functions */

var next = 0

func (chord *ChordNode) notifed(newPredecessor *Node) {
	if debugging {
		fmt.Println("predecessor:", newPredecessor)
	}
	if chord.Predecessor == nil ||
		between1(chord.Predecessor.Id, newPredecessor.Id,
			chord.LocalNode.Id, false) {
		chord.Predecessor = newPredecessor
		return
	}
	if debugging {
		fmt.Println("Fuck Off <Don't need a new predecessor>")
	}
}

func (chord *ChordNode) stabilize() {
	if debugging {
		fmt.Println("<Stabilizing>")
	}
	list := make([]*Node, m)
	var successors []*Node
	var predecessor *Node

	for {
		var err error
		var err1 error
		predecessor, err = chord.getPredecessor()
		successors, err1 = chord.getSuccessors()
		if err != nil || err1 != nil {
			if debugging {
				fmt.Println("Failed to contact Successor:", err)
				fmt.Println(err1)
			}
			for i := 0; i < m-1; i++ {
				chord.Successor[i] = chord.Successor[i+1]
			}
			chord.Successor[m-1] = chord.LocalNode
			return
		}
		break
	}

	if predecessor != nil {
		if between1(chord.LocalNode.Id, predecessor.Id, chord.Successor[0].Id, false) {
			chord.Successor[0] = predecessor
		}
	}
	list[0] = chord.Successor[0]
	for i := 1; i < m; i++ {
		list[i] = successors[i-1]
	}
	chord.Successor = list
	chord.notify()

}

func (chord *ChordNode) check_predecessor() {
	if chord.Predecessor != nil {
		if !chord.isAlive() {
			chord.Predecessor = nil
		}
	}
}

func (chord *ChordNode) fix_fingers() {
	node := chord.find(chord.LocalNode, chord.FingerTable[next].Start)
	chord.FingerTable[next].Successor = node
	next = next + 1

	if next >= keySize {
		next = 0
	}

	for {
		previous := next - 1
		if previous < 0 {
			previous = keySize - 1
		}

		if between1(chord.FingerTable[previous].Start,
			chord.FingerTable[next].Start, node.Id, true) {

			chord.FingerTable[next].Successor = node
			next = next + 1
			if next >= keySize {
				next = 0
			}
			continue
		}
		break
	}
}

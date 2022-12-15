package main

import "fmt"

/* Finding Successor Functions */

func (chord *ChordNode) find_successor(id string) (bool, *Node) {
	if chord.Predecessor != nil &&
		between1(chord.Predecessor.Id, id, chord.LocalNode.Id, true) {
		return true, chord.LocalNode
	}

	if between1(chord.LocalNode.Id, id, chord.Successor[0].Id, true) {
		return true, chord.Successor[0]
	} else {
		return chord.closest_preceding_node(id)
		//return false, chord.Successor[0]
	}
}

func (chord *ChordNode) closest_preceding_node(id string) (bool, *Node) {
	for i := keySize - 1; i >= 0; i-- {
		if between1(chord.LocalNode.Id, chord.FingerTable[i].Start,
			id, false) {
			return false, chord.FingerTable[i].Successor
		}
	}

	return false, chord.LocalNode
}

// Search iteratively, if not found asking the received node next
func (chord *ChordNode) find(startNode *Node, id string) *Node {
	maxSteps := keySize - 1
	var n = startNode
	var succesor *Node
	var flag bool

	for i := 0; i < maxSteps; i++ {
		flag, succesor = chord.find_succesor(n, id)
		if flag == true {
			return succesor
		}
		n = succesor
	}
	if debugging {
		fmt.Println("Couldnt find Succesor, returning start node")
	}
	return startNode
}

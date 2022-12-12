package main

/* Finding Successor Functions */

func (chord *ChordNode) find_successor(id string) (bool, *Node) {
	for {
		if between1(chord.LocalNode.Id, id, chord.Successor[0].Id, true) {
			return true, chord.Successor[0]
		} else {
			return chord.closest_preceding_node(id)
			//return false, chord.Successor[0]
		}
	}
}

func (chord *ChordNode) closest_preceding_node(id string) (bool, *Node) {
	for i := 159; i >= 0; i-- {
		if between1(chord.LocalNode.Id, chord.FingerTable[i].Start,
			id, true) {
			return false, chord.FingerTable[i].Successor
		}
	}
	//return chord.Successor[0]
	//return true, chord.Successor[0]
	return false, chord.Successor[0]
}

func (chord *ChordNode) find(startNode *Node, id string) *Node {
	maxSteps := 20
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

	return startNode
}

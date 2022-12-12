package main

import (
	"fmt"
	"net/rpc"
	"strconv"
)

/* RPC Calling Functions */

func (chord *ChordNode) call(address string, method string, request RpcArgs) (RpcArgs, error) {
	var reply RpcArgs
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		if debugging {
			fmt.Println("RPC dial error:", err)
		}
		return reply, err
	}
	defer client.Close()
	if err := client.Call(method, request, &reply); err != nil {
		if debugging {
			fmt.Println("Client call error:", err)
		}
		return reply, err
	}
	if debugging {
		fmt.Println(reply.Value)
	}
	return reply, nil
}

func (chord *ChordNode) getPredecessor() (*Node, error) {
	address := chord.Successor[0].Address
	response, err := chord.call(address, "ChordNode.SendPredecessor", RpcArgs{"0", "", nil, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to get predecessor!!!!!! <", err)
		}
		return nil, err
	}
	return response.RNode, nil
}

func (chord *ChordNode) getSuccessors() ([]*Node, error) {
	address := chord.Successor[0].Address
	response, err := chord.call(address, "ChordNode.SendSuccessors", RpcArgs{})
	if err != nil {
		if debugging {
			fmt.Println("Client call error:", err)
			return nil, err
		}
	}
	if debugging {
		fmt.Println("Sent Successors")
	}
	if len(response.RNodes) < m {
		for i := len(response.RNodes); i < m; i++ {
			response.RNodes = append(response.RNodes, chord.LocalNode)
		}
	}

	if len(response.RNodes) > m {
		response.RNodes = response.RNodes[:m]
	}

	return response.RNodes, nil
}

func (chord *ChordNode) notify() {
	address := chord.Successor[0].Address
	_, err := chord.call(address, "ChordNode.Notified", RpcArgs{"", "Notifying", chord.LocalNode, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to send Notification! <", err)
		}
	}
}

func (chord *ChordNode) find_succesor(sendTo *Node, id string) (bool, *Node) {
	address := sendTo.Address
	reply, err := chord.call(address, "ChordNode.FindingSuccesor", RpcArgs{id, "Find The Successor", nil, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to find_succesor request! <", err)
		}
		return false, chord.Successor[0]
	}

	flag, _ := strconv.ParseBool(reply.Value)
	node := reply.RNode

	return flag, node
}

func (chord *ChordNode) put_all() {
	address := chord.Successor[0].Address
	_, err := chord.call(address, "ChordNode.Put_all", RpcArgs{"", "", nil, nil, chord.Bucket})
	if err != nil {
		if debugging {
			fmt.Println("Failure to hand over keys! <", err)
		}
	}
	for _, value := range chord.Bucket {
		postSender(address, value) // Send stored files to new host
	}
}

func (chord *ChordNode) get_all() map[string]string {
	address := chord.Successor[0].Address
	reply, err := chord.call(address, "ChordNode.Get_all", RpcArgs{"", "GET", chord.LocalNode, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to send Notification! <", err)
		}
	}
	return reply.Keys
}

func (chord *ChordNode) isAlive() bool {
	address := chord.Predecessor.Address
	reply, err := chord.call(address, "ChordNode.Alive", RpcArgs{"", "Alive?", nil, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Predecessor Down!!:", err)
		}
		return false
	}
	if debugging {
		fmt.Println(reply.Key + ": " + reply.Value)
	}
	return true
}

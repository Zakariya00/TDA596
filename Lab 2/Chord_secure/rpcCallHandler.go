package main

import (
	"fmt"
	"net/rpc"
	"strconv"
)

/* RPC Calling Functions */

// call helper Function, takes care of dialing and calling client and returns result
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

// getPredecessor Call to ask for predecessor
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

// getSuccessors Call to ask for successors
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

// notify successor
func (chord *ChordNode) notify() {
	address := chord.Successor[0].Address
	_, err := chord.call(address, "ChordNode.Notified", RpcArgs{"", "Notifying", chord.LocalNode, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to send Notification! <", err)
			return
		}
	}
}

// find_succesor Call to "find" successor
func (chord *ChordNode) find_succesor(sendTo *Node, id string) (bool, *Node) {
	address := sendTo.Address
	reply, err := chord.call(address, "ChordNode.FindingSuccesor", RpcArgs{id, "Find The Successor", nil, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failure to find_succesor request! <", err)
		}
		return false, sendTo
	}

	flag, _ := strconv.ParseBool(reply.Value)
	node := reply.RNode

	return flag, node
}

// put_all Call to hand over keys before shutdown
func (chord *ChordNode) put_all() {
	if len(chord.Bucket) == 0 && len(chord.Backups) == 0 {
		fmt.Println("No key/s to hand over")
		return
	}
	for i := 0; i < m; i++ {
		if chord.Successor[i].Id == chord.LocalNode.Id {
			continue
		}
		address := chord.Successor[i].Address
		_, err := chord.call(address, "ChordNode.Put_all", RpcArgs{"", "", nil, nil, chord.Bucket})
		if err != nil {
			fmt.Printf("Sending keys to succesor[%d] failed\n", i)
			continue
		}

		for _, value := range chord.Bucket {
			postSender(address, value)
		}
		for key, value := range chord.Backups {
			chord.sendBackup(key, value)
		}

		fmt.Printf("Sent keys to succesor[%d] ---> %+v\n", i, *chord.Successor[i])
		return
	}
	fmt.Println("Failed to hand over keys to successor/s")
}

// get_all Call to ask successor for your keys
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

// isAlive Call to check if predecessor is still alive
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

// sendBackup sends backup to successor
func (chord *ChordNode) sendBackup(key, value string) error {
	address := chord.Successor[0].Address
	reply, err := chord.call(address, "ChordNode.GetBackUp", RpcArgs{key, value, nil, nil, nil})
	if err != nil {
		if debugging {
			fmt.Println("Failed to send backup!!:", err)
		}
		return err
	}
	if debugging {
		fmt.Println(reply.Value)
	}
	postSender(address, value)
	return nil
}

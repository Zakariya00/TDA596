package main

import (
	"fmt"
	"strconv"
)

/* RPC Receiving Handler Functions */

// Ping Connectivity check, for testing
func (chord *ChordNode) Ping(args RpcArgs, reply *RpcArgs) error {
	fmt.Println("<Ping>:", args.Value)
	*reply = RpcArgs{"", "Ping Back", nil, nil, nil}
	return nil
}

// Get Check for requested key, send back the value
func (chord *ChordNode) Get(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Get>")
	}
	value := chord.Bucket[args.Key]
	if debugging {
		fmt.Println("Lookup for Key:", args.Key)
		fmt.Println("Found value:", value)
	}
	*reply = RpcArgs{"", value, nil, nil, nil}
	return nil
}

// Put the received key & value in bucket
func (chord *ChordNode) Put(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Put>")
	}
	chord.Bucket[args.Key] = args.Value
	if debugging {
		fmt.Println("Stored Key & Value pair:", args.Key, args.Value)
	}
	*reply = RpcArgs{"", "Stored", nil, nil, nil}
	return nil
}

// Delete the received key & value from bucket
func (chord *ChordNode) Delete(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Delete>")
	}
	delete(chord.Bucket, args.Key)
	if debugging {
		fmt.Println("Deleted key:", args.Key)
	}
	*reply = RpcArgs{"", "Deleted", nil, nil, nil}
	return nil
}

// SendPredecessor send back predecessor
func (chord *ChordNode) SendPredecessor(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Sending Predecessor>")
	}
	*reply = RpcArgs{"", "Sent Predecessor", chord.Predecessor, nil, nil}
	return nil
}

// SendSuccessors send back successors
func (chord *ChordNode) SendSuccessors(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Sending Successors>")
	}
	*reply = RpcArgs{"", "", nil, chord.Successor, nil}
	return nil
}

// Notified check then/if set new predecessor
func (chord *ChordNode) Notified(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Notifying Node>")
	}
	chord.notified(args.RNode)
	*reply = RpcArgs{"", "Node Notified", nil, nil, nil}
	return nil
}

// FindingSuccesor search for the successor to received key, return results
func (chord *ChordNode) FindingSuccesor(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Finding Succesor>")
	}
	flag, node := chord.find_successor(args.Key)
	*reply = RpcArgs{"", strconv.FormatBool(flag), node, nil, nil}
	return nil
}

// Put_all put all received keys in your bucket
func (chord *ChordNode) Put_all(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Putting All>")
	}
	received := args.Keys
	for key, value := range received {
		chord.Bucket[key] = value
	}
	*reply = RpcArgs{"", "Received Bucket", nil, nil, nil}
	return nil
}

// Get_all send back keys that belong to your new predecessor
func (chord *ChordNode) Get_all(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Getting All>")
	}
	newBucket := make(map[string]string)
	if len(chord.Bucket) != 0 {
		for key, value := range chord.Bucket {
			if between1(chord.Predecessor.Id, key, args.RNode.Id, true) {
				newBucket[key] = value
				delete(chord.Bucket, key)
				go postSender(args.RNode.Address, value) // send file to new host
			}
		}
	}
	*reply = RpcArgs{"", "", nil, nil, newBucket}
	return nil
}

// Alive send back reply to show node is still running
func (chord *ChordNode) Alive(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println(args.Value)
	}
	*reply = RpcArgs{args.Value, "Yeah Im alive", nil, nil, nil}
	return nil
}

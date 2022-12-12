package main

import (
	"fmt"
	"strconv"
)

/* RPC Receiving Handler Functions */

func (chord *ChordNode) Ping(args RpcArgs, reply *RpcArgs) error {
	fmt.Println("<Ping>:", args.Value)
	*reply = RpcArgs{"", "Ping Back", nil, nil, nil}
	return nil
}

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

func (chord *ChordNode) SendPredecessor(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Sending Predecessor>")
	}
	*reply = RpcArgs{"", "Sent Predecessor", chord.Predecessor, nil, nil}
	return nil
}

func (chord *ChordNode) SendSuccessors(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Sending Successors>")
	}
	*reply = RpcArgs{"", "", nil, chord.Successor, nil}
	return nil
}

func (chord *ChordNode) Notified(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Notifying Node>")
	}
	chord.notifed(args.RNode)
	*reply = RpcArgs{"", "Node Notified", nil, nil, nil}
	return nil
}

func (chord *ChordNode) FindingSuccesor(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println("<Finding Succesor>")
	}
	flag, node := chord.find_successor(args.Key)
	*reply = RpcArgs{"", strconv.FormatBool(flag), node, nil, nil}
	return nil
}

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

func (chord *ChordNode) Alive(args RpcArgs, reply *RpcArgs) error {
	if debugging {
		fmt.Println(args.Value)
	}
	*reply = RpcArgs{args.Value, "Yeah Im alive", nil, nil, nil}
	return nil
}

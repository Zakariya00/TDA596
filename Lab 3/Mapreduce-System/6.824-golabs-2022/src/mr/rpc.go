package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"time"
)
import "strconv"

// Add your RPC definitions here.

type TypeOfTask int

const (
	MAP = iota
	REDUCE
	SUSPEND
)

type StateOfTask int

const (
	FREE = iota
	TAKEN
	FINISHED
	TIMEOUT
)

type Task struct {
	index int         // Task index
	file  string      // Task file
	start time.Time   // Task start time
	state StateOfTask // Task state
	taskT TypeOfTask  // Task type

}

// Task completed sending args
type Targs struct {
	Index int        // Completed Task index
	File  string     // Completed Task file
	Type  TypeOfTask // Completed Task type
}

// Task request reply args
type Treply struct {
	Index  int        // Assigned Task index
	File   string     // Assigned Task file
	Type   TypeOfTask // Assigned Task type
	Reduce int        //
}

type Tstatus struct {
	Status bool   // Reported Task status (true = successfully completed & approved)
	Msg    string // Reported Task message
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

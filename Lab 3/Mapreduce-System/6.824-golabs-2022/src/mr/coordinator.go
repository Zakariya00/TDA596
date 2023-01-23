package mr

import (
	"log"
	"sync"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	// Your definitions here.
	lock   sync.Mutex // Lock for synchronization
	mTasks []Task     // Map tasks
	rTasks []Task     // Reduce tasks
	mRem   int        // Map tasks remaining
	rRem   int        // Reduce tasks remaining
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) SendTask(args *Targs, reply *Treply) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If there are remaining Map tasks, assigns them to workers
	if c.mRem > 0 {
		for i, e := range c.mTasks {
			if e.state != FINISHED && e.state != TAKEN {
				*reply = Treply{e.index, e.file, e.taskT, len(c.rTasks)}
				c.mTasks[i].state = TAKEN
				c.mTasks[i].start = time.Now()
				return nil
			}
		}

		// If there are no more remaining Map tasks, assigns Reduce task to workers
	} else if c.rRem > 0 {
		for i, e := range c.rTasks {
			if e.state != FINISHED && e.state != TAKEN {
				*reply = Treply{e.index, e.file, e.taskT, len(c.rTasks)}
				c.rTasks[i].state = TAKEN
				c.rTasks[i].start = time.Now()
				return nil
			}
		}
		return nil
	}

	// If there are no remaining Map or Reduce tasks, Suspends worker
	*reply = Treply{-1, "", SUSPEND, -1}
	return nil
}

// Receives and Handles, reported finished Task
func (c *Coordinator) TaskDone(args *Targs, reply *Tstatus) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var status = false
	var msg string

	// Checks completed Task type, then marks that task as finished
	switch args.Type {
	case MAP:
		// Check if Task has not already been Completed
		if c.mTasks[args.Index].state != FINISHED {
			c.mTasks[args.Index].state = FINISHED
			c.mRem--

			status = true
			msg = "Map Task Successfully Completed"
		}

	case REDUCE:
		// Check if Task has not already been Completed
		if c.rTasks[args.Index].state != FINISHED {
			c.rTasks[args.Index].state = FINISHED
			c.rRem--

			status = true
			msg = "Reduce Task Successfully Completed"
		}
	}

	*reply = Tstatus{status, msg}
	return nil
}

// Checks for Task timeouts and marks them
func (c *Coordinator) timeoutHandler() {
	defer c.lock.Unlock()
	for {
		c.lock.Lock()
		if c.mRem > 0 {
			for i, e := range c.mTasks {
				// Checks for timeout
				if e.start.Before(time.Now().Add(-10 * time.Second)) {
					c.mTasks[i].state = TIMEOUT
				}
			}
		} else if c.rRem > 0 {
			for i, e := range c.rTasks {
				// Checks for timeout
				if e.start.Before(time.Now().Add(-10 * time.Second)) {
					c.rTasks[i].state = TIMEOUT
				}
			}
		}
		c.lock.Unlock()
	}
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {

	// Your code here.
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.mRem == 0 && c.rRem == 0
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		mTasks: make([]Task, len(files)),
		rTasks: make([]Task, nReduce),
		mRem:   len(files),
		rRem:   nReduce,
	}

	// Initialize Map Tasks
	for i, file := range files {
		c.mTasks[i] = Task{index: i, file: file, state: FREE, taskT: MAP}
	}

	// Initialize Reduce Tasks
	for i := 0; i < nReduce; i++ {
		c.mTasks[i] = Task{index: i, state: FREE, taskT: REDUCE}
	}

	go c.timeoutHandler()

	c.server()
	return &c
}

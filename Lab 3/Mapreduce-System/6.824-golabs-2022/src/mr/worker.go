package mr

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)
import "log"
import "net/rpc"
import "hash/fnv"

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	for {
		task := requestTask()
		switch task.Type {
		case MAP:
			mapHandler(task, mapf)

		case REDUCE:
			reduceHandler(task, reducef)

		case SUSPEND:
			time.Sleep(1000 * time.Millisecond)

		}
	}
}

// Handle Map Task
func mapHandler(task Treply, mapf func(string, string) []KeyValue) {
	// Open & read file content
	file, err := os.Open(task.File)
	if err != nil {
		log.Fatal("File open failure <%s>", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("File read failure <%s>", err)
	}

	// Map function
	intermediate := mapf(task.File, string(content))
	rFiles := make(map[int][]KeyValue)

}

// Handle Reduce Task
func reduceHandler(task Treply, reducef func(string, []string) string) {

}

// Request New Task
func requestTask() Treply {
	args := Targs{}
	reply := Treply{}
	ok := call("Coordinator.SendTask", &args, &reply)
	if !ok {
		log.Fatal("Task request call failure")
	}
	return reply
}

// Request Finished Task
func reportTask(task Treply) Tstatus {
	args := Targs{
		Index: task.Index,
		File:  task.File,
		Type:  task.Type,
	}
	reply := Tstatus{}
	ok := call("Coordinator.TaskDone", &args, &reply)
	if !ok {
		log.Fatal("Task report done call failure")
	}
	return reply
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

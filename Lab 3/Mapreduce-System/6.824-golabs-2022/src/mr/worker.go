package mr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
	var task = Treply{}
	var stop = false
	for !stop {
		task, stop = requestTask()
		switch task.Type {

		case MAP:
			mapHandler(task, mapf)
			_, stop = reportTask(task)

		case REDUCE:
			reduceHandler(task, reducef)
			_, stop = reportTask(task)

		case SUSPEND:
			time.Sleep(1000 * time.Millisecond)

		default:
			continue
		}

	}
}

// Handle Map Task
func mapHandler(task Treply, mapf func(string, string) []KeyValue) {
	// Open file
	file, err := os.Open(task.File)
	if err != nil {
		log.Fatal("File open failure <%s>", err)
	}
	defer file.Close()

	// Read file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("File read failure <%s>", err)
	}

	// Map function
	kva := mapf(task.File, string(content))

	// Gather intermediate keys with the same hash
	intermediate := make(map[int][]KeyValue)
	for _, kv := range kva {
		hashID := ihash(kv.Key) % task.Reduce
		intermediate[hashID] = append(intermediate[hashID], kv)
	}

	// Store gathered intermediate keys in files
	for rTaskID, kvs := range intermediate {
		filename := fmt.Sprintf("mr-%d-%d", task.Index, rTaskID)
		tmpfilename := fmt.Sprintf("mr-%d-%d-%d", task.Index, rTaskID, os.Getpid())
		file, err := os.Create(tmpfilename)
		if err != nil {
			log.Fatal("Intermediate File create failure <%s>", err)
		}
		defer file.Close()

		// Write in JSON format to file
		enc := json.NewEncoder(file)
		for _, kv := range kvs {
			err := enc.Encode(&kv)
			if err != nil {
				log.Fatal("Intermediate File encoding failure <%s>", err)
			}
		}
		// Rename file
		os.Rename(tmpfilename, filename)
	}
}

// Handle Reduce Task
func reduceHandler(task Treply, reducef func(string, []string) string) {
	// Get all files for reduce
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Get current Directory failure <%s>", err)
	}
	files, err := filepath.Glob(fmt.Sprintf("%v/mr-%d-%d", currDir, "*", task.Index))
	if err != nil {
		log.Fatal("Geting reduce files failure <%s>", err)
	}

	intermediate := []KeyValue{}
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("Reduce file open failure <%s>", err)
		}
		defer file.Close()

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
	}

	// -----------------------------------------------------------
	sort.Sort(ByKey(intermediate))
	oname := fmt.Sprintf("mr-out-%d", task.Index)
	tmponame := fmt.Sprintf("mr-out-%d-%d", task.Index, os.Getpid())
	ofile, _ := os.Create(tmponame)
	defer ofile.Close()

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	os.Rename(tmponame, oname)
}

// Request New Task
func requestTask() (Treply, bool) {
	args := Targs{}
	reply := Treply{}
	ok := call("Coordinator.SendTask", &args, &reply)
	if !ok {
		// Cant reach Coordinator, Assume no More Work and Exit
		return Treply{}, true
	}
	return reply, false
}

// Request Finished Task, handle Exit too
func reportTask(task Treply) (Tstatus, bool) {
	args := Targs{
		Index: task.Index,
		File:  task.File,
		Type:  task.Type,
	}
	reply := Tstatus{}
	ok := call("Coordinator.TaskDone", &args, &reply)
	if !ok {
		// Cant reach Coordinator, Assume no More Work and Exit
		return Tstatus{}, true
	}
	return reply, false
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

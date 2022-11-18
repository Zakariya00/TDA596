package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var port string
var nbrOfClients = 0
var lock sync.Mutex

func main() {
	port = os.Args[1]
	if len(port) == 0 {
		port = ":80"
	}
	fmt.Println("Port:", port)

	//Establish a socket connection
	listener, error := net.Listen("tcp", port)
	if error != nil {
		fmt.Println("Error on server start")
		listener.Close()
		os.Exit(1)
	}

	for {
		//Wait for incoming client connections
		//Each new client request is accepted
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println("Error on client connection")
			connection.Close()
			continue
		}

		//To avoid overwhelming your server, you should not create more than a reasonable number of child processes
		//(for this assignment, use at most 10), in which case your server should wait until one of its ongoing child
		//processes exits before forking a new one to handle the new request.
		for nbrOfClients >= 10 {
		}

		//Pass client connection to spawned child process to handle Client
		go handleClient(connection)

	}
}

// Handles client entrance and exit
func handleClient(connection net.Conn) {
	lock.Lock()
	nbrOfClients++
	lock.Unlock()

	reader := bufio.NewReader(connection)
	req, error := http.ReadRequest(reader)

	//res, _ := httputil.DumpRequest(req, true)
	if error != nil {
		fmt.Println("Bad request")
		errorHandler(connection, 404, "Unable To Serve Bad Request")
		return
	}
	//fmt.Println(res)

	handleRequest(connection, req)
	connection.Close()

	lock.Lock()
	nbrOfClients--
	lock.Unlock()
}

// Handles client request
func handleRequest(connection net.Conn, req *http.Request) {
	fmt.Println(req.Method)

	switch req.Method {
	case "GET":
		getHandler(connection, req)
	case "POST":
		postHandler(connection, req)
	default:
		errorHandler(connection, 501, "Method Not Implemented")
	}
}

// Handles GET requests
func getHandler(connection net.Conn, req *http.Request) {
	//Ignore annoying GET requests for favicon.ico
	if path.Base(req.URL.Path) == "favicon.ico" {
		return
	}

	fullPath := req.URL
	filename := path.Base(req.URL.Path)
	contentType := req.Header.Get("Content-type")

	fmt.Println(fullPath)
	fmt.Println(filename)
	fmt.Println(contentType)

	//fileType := req.Header.Get("Content-Type")
	var extension = filepath.Ext(filename)
	var fileType string
	if len(extension) > 0 {
		fileType = extension[1:len(extension)]
	} else {
		fileType = extension
	}

	if fullPath.String() == "/" {

	} else {

		switch fileType {
		//Server should accept requests for files ending in
		//html, txt, gif, jpeg, jpg, or css
		case "html", "txt", "gif", "jpeg", "jpg", "css":
			fmt.Println(fileType)
		default:
			//If the client requests a file with any other extension,
			//the web server must respond with a well-formed 400 "Bad Request" code
			errorHandler(connection, 400, "")
			return
		}
		fileGetHandler(connection, filename, fileType)
		return
	}

	//Transmit them to the client with a Content-Type of
	//text/html, text/plain, image/gif, image/jpeg, image/jpeg,
	//or text/css, respectively
	body := "Hello Zakariya, Welcome"
	fmt.Fprint(connection, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(connection, "Content-Type: text/html\r\n")
	fmt.Fprint(connection, "\r\n")
	fmt.Fprint(connection, body)
}

func fileGetHandler(connection net.Conn, fPath string, ftype string) {
	fmt.Println("Sending file")
	myDir, error := os.Getwd()
	if error != nil {
		fmt.Println("Error on finding working directory")
		fmt.Println(error)
	}
	fmt.Println(myDir)

	info, err := os.Stat(fPath)
	if os.IsNotExist(err) {
		fmt.Println("File doesnt exist")
		errorHandler(connection, 404, "File doesnt exists")
		return
	}
	size := info.Size()

	//file, error1 := os.Open(fPath)

	file, error1 := os.Open(fPath)
	if os.IsNotExist(error1) {
		fmt.Println("writing File failed")
		return
	}

	//file.Read()
	writer := bufio.NewWriter(connection)

	fmt.Fprint(connection, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(connection, "Content-Length: %d\r\n", size)
	fmt.Fprintf(connection, "Content-Type: %s\r\n", contentTypeHandler(ftype))
	fmt.Fprint(connection, "\r\n")
	fmt.Fprint(connection, file)
	io.Copy(writer, file)
}

// Handles POST requests
func postHandler(connection net.Conn, req *http.Request) {
	// For POST requests, please make sure that you store the files
	//and make them accessible with a subsequent  GET request

	fmt.Println(req.Method)
	fmt.Println(req.Body)

	f, err := os.Create(path.Base(req.URL.Path))
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, req.Body)
	if err != nil {
		fmt.Println("Error")
		return
	}

	fmt.Fprint(connection, "HTTP/1.1 201 CREATED\r\n")
}

func contentTypeHandler(ftype string) string {
	switch ftype {
	case "html":
		return "text/html"
	case "txt":
		return "text/plain"
	case "gif":
		return "image/gif"
	case "jpeg":
		return "image/jpeg"
	case "jpg":
		return "image/jpeg"
	case "css":
		return "text/css"
	}

	return ""
}

// Handles any faulty or unsupported requests
func errorHandler(connection net.Conn, errorCode int, msg string) {
	if errorCode == 501 {
		fmt.Fprint(connection, "HTTP/1.1 501 Not Implemented\r\n")
	} else if errorCode == 404 {
		fmt.Fprint(connection, "HTTP/1.1 404 Not Found\r\n")
	} else {
		fmt.Fprint(connection, "HTTP/1.1 400 Bad Request\r\n")
	}

	body := msg

	fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(connection, "Content-Type: text/html\r\n")
	fmt.Fprint(connection, "\r\n")
	fmt.Fprint(connection, body)
}

func sendHandler(connection net.Conn, body string, header string, contentType string) {
	fmt.Fprint(connection, "%s\r\n", header)
	fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(connection, "Content-Type: %s\r\n", contentType)
	fmt.Fprint(connection, "\r\n")
	fmt.Fprint(connection, body)
}

/*
curl -i -X GET http://localhost:8080
*/

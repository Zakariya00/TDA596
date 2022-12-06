package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var port string

const maxClients = 10 // Maximum client that server handles simultaneously
var sem chan int      // Semaphore, syncing mechanism for Server

func main() {
	sem = make(chan int, maxClients) // Initialize semaphore for maxClient

	//set port as command arg or else default port (8080)
	if len(os.Args) < 2 {
		port = ":8080"
	} else {
		port = os.Args[1]
	}
	fmt.Printf("Port%v\n", port)

	//Establish a socket connection
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error on server start: %v\n", err)
		listener.Close()
		os.Exit(1)
	}
	defer listener.Close()

	//Handlers
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	http.HandleFunc("/", serverHandler)

	//Listens for and serves http request by spawning a go process with handler
	if err = http.Serve(listener, nil); err != nil {
		fmt.Printf("Error serving HTTP requests: %v\n", err)
		os.Exit(1)
	}
}

// Main handler for http requests to server, with sync mechanism
func serverHandler(w http.ResponseWriter, r *http.Request) {
	sem <- 1
	fmt.Println("-------------- <Sem Acquired> -------------- ")
	defer func() {
		<-sem
		fmt.Println("-------------- <Sem Released> -------------- ")
	}()

	//Print for Debugging
	fmt.Printf("File-Path: %s\n", r.URL)
	fmt.Printf("File-Name: %s\n", path.Base(r.URL.Path))
	fmt.Printf("Request-Method: %s\n", r.Method)

	if !checkMethod(w, r) {
		return
	}
	if r.Method == "POST" {
		postHandler(w, r)
	} else {
		getHandler(w, r, getFileType(path.Base(r.URL.Path)))
	}
}

// Handles GET requests
func getHandler(w http.ResponseWriter, r *http.Request, fType string) {
	if path.Base(r.URL.Path) == "/" {
		fmt.Fprintf(w, "Hello, Welcome to the Main Page")
		return
	} else if !validFileType(w, fType) {
		return
	}

	http.ServeFile(w, r, "./"+r.URL.String())
}

// Handles POST requests
func postHandler(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error on retrieving file <POST>: "+err.Error(), 500)
		return
	}
	defer file.Close()

	//check for valid file type
	if !validFileType(w, getFileType(handler.Filename)) {
		return
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create Local file
	f, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, "Error on create file <POST>: "+err.Error(), 500)
		return
	}
	defer f.Close()

	// Write to local file within our directory that follows
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error on writing to file <POST>: "+err.Error(), 500)
		return
	}

	// return that we have successfully uploaded client file!
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully Uploaded File,\nFeel free to visit again!!"))
}

// Check for valid Http Requesting Method
func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	switch r.Method {
	case "GET":
		return true
	case "POST":
		return true
	default:
		http.Error(w, "Request Method Is Currently Not Supported <"+r.Method+">", 501)
		return false
	}
}

// Get file type from request file name
func getFileType(filename string) string {
	var fType string
	var extension = filepath.Ext(filename)
	if len(extension) > 0 {
		fType = extension[1:]
	} else {
		fType = extension
	}
	return fType
}

// Check for a valid file type
func validFileType(w http.ResponseWriter, fType string) bool {
	switch fType {
	case "html", "txt", "gif", "jpeg", "jpg", "css":
		return true
	default:
		http.Error(w, "File type not supported", 400)
		return false
	}
}

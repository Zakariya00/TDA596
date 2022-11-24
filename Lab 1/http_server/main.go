package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
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
		postHandler(w, r, getFileType(r))
	} else {
		getHandler(w, r, getFileType(r))
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
func postHandler(w http.ResponseWriter, r *http.Request, fType string) {
	if path.Base(r.URL.Path) == "/" {
		http.Error(w, "Need a file name in the URL for a valid POST request", 400)
		return
	} else if !validFileType(w, fType) {
		return
	}

	defer r.Body.Close()
	//Create file with the given file name in URL Request
	f, err := os.Create(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "Error on create file <POST>: "+err.Error(), 500)
		return
	}
	defer f.Close()

	//Process sent file by file type
	switch fType {
	case "jpeg", "jpg", "gif":
		{
			img, _, err := image.Decode(r.Body)
			if err != nil {
				http.Error(w, "Error on decoding request data to image <POST>: "+err.Error(), 500)
				return
			}

			if fType == "gif" {
				err = gif.Encode(f, img, nil)
			} else {
				err = jpeg.Encode(f, img, nil)
			}

			if err != nil {
				http.Error(w, "Error on writing image to file <POST>: "+err.Error(), 500)
				return
			}
		}
	default:
		{
			_, err = io.Copy(f, r.Body)
			if err != nil {
				http.Error(w, "Error on writing to file <POST>: "+err.Error(), 500)
				return
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
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

// Get file type from request URL
func getFileType(r *http.Request) string {
	filename := path.Base(r.URL.Path)
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

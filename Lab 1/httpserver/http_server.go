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

const maxClients = 10

// var semaphore chan
var sem chan int

func main() {
	sem = make(chan int, maxClients)

	if len(os.Args) < 2 {
		port = ":8080"
	} else {
		port = os.Args[1]
	}
	fmt.Println("Port", port)

	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldnt find working directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Home directory Path: %s\n", path)

	//Establish a socket connection
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error on server start: %v\n", err)
		listener.Close()
		os.Exit(1)
	}

	//Handlers
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", rootHandler)

	//semaphore <- struct{}{}
	//Serve http requests that come through
	if err = http.Serve(listener, nil); err != nil {
		fmt.Printf("Error serving HTTP requests: %v\n", err)
		//	func() { <-semaphore }()
	}
}

// Added to get rid of the annoying requests
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	func() {
		sem <- 1
		fmt.Println("-------------- <Sem Acquired> -------------- ")
	}()
	defer func() {
		<-sem
		fmt.Println("-------------- <Sem Released> -------------- ")
	}()
	http.ServeFile(w, r, "relative/path/to/favicon.ico")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	func() {
		sem <- 1
		fmt.Println("-------------- <Sem Acquired> -------------- ")
	}()
	defer func() {
		<-sem
		fmt.Println("-------------- <Sem Released> -------------- ")
	}()

	fmt.Printf("File-path: %s\n", r.URL)
	fmt.Printf("File-name: %s\n", path.Base(r.URL.Path))

	if !checkMethod(w, r) {
		return
	}

	if r.Method == "POST" {
		postHandler(w, r, getFileType(r))
	} else {
		getHandler(w, r, getFileType(r))
	}
}

func getHandler(w http.ResponseWriter, r *http.Request, fType string) {
	if path.Base(r.URL.Path) == "/" {
		fmt.Fprintf(w, "Hello, Welcome to the Main Page")
		return
	} else if !validFileType(w, fType) {
		return
	}

	fmt.Println("---------- <sending file> ----------")
	http.ServeFile(w, r, "./"+r.URL.String())
}

// Works except for corruption of photofiles
func postHandler(w http.ResponseWriter, r *http.Request, fType string) {
	if path.Base(r.URL.Path) == "/" {
		http.Error(w, "Need a file name in the URL for a valid POST request", 400)
		return
	} else if !validFileType(w, fType) {
		return
	}

	fmt.Println("---------- <receiving file> ----------")
	defer r.Body.Close()

	f, err := os.Create(path.Base(r.URL.Path))
	if err != nil {
		fmt.Printf("Error on create file <POST>: %v\n", err)
		http.Error(w, "Error on create file <POST>: "+err.Error(), 500)

		return
	}

	switch fType {
	case "jpeg", "jpg", "gif":
		{
			img, _, err := image.Decode(r.Body)
			if err != nil {
				fmt.Printf("Error on decoding request data to image <POST>: %v\n", err)
				http.Error(w, "Error on decoding request data to image <POST>: "+err.Error(), 500)
				return
			}

			if fType == "gif" {
				err = gif.Encode(f, img, nil)
			} else {
				err = jpeg.Encode(f, img, nil)
			}

			if err != nil {
				fmt.Printf("Error on writing image to file <POST>: %v\n", err)
				http.Error(w, "Error on writing image to file <POST>: "+err.Error(), 500)
				return
			}

		}
	default:
		{
			_, err = io.Copy(f, r.Body)
			if err != nil {
				fmt.Printf("Error on writing to file <POST>: %v\n", err)
				http.Error(w, "Error on writing to file <POST>: "+err.Error(), 500)
				return
			}
		}
	}

	fmt.Println("---------- <received file> ----------")
	w.WriteHeader(http.StatusCreated)
}

// Helper Methods
func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	fmt.Printf("Checked Method: %s\n", r.Method)

	switch r.Method {
	case "GET":
		return true
	case "POST":
		return true
	default:
		http.Error(w, "Request Method Is Currently Not Supported", 501)
		return false
	}
}

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

func validFileType(w http.ResponseWriter, fType string) bool {
	switch fType {
	case "html", "txt", "gif", "jpeg", "jpg", "css":
		return true
	default:
		http.Error(w, "File type not supported", 400)
		return false
	}
}

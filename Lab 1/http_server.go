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

func main() {
	if len(os.Args) < 2 {
		port = ":8080"
	} else {
		port = os.Args[1]
	}
	fmt.Println("Port", port)

	//Establish a socket connection
	listener, error := net.Listen("tcp", port)
	if error != nil {
		fmt.Printf("Error on server start: %v\n", error)
		listener.Close()
		os.Exit(1)
	}

	path, err := os.Getwd()
	if error != nil {
		fmt.Printf("Couldnt find working directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Home directory Path: %s\n", path)

	//Handlers
	http.HandleFunc("/", rootHandler)
	//http.Handle("/", fileHandler)

	//Serve http requests that come through
	if error = http.Serve(listener, nil); error != nil {
		fmt.Printf("Error serving HTTP requests: %v\n", error)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL
	filename := path.Base(r.URL.Path)
	var fType string

	var extension = filepath.Ext(filename)
	if len(extension) > 0 {
		fType = extension[1:len(extension)]
	} else {
		fType = extension
	}

	if filename == "favicon.ico" || !checkMethod(w, r) {
		return
	}

	if filename == "/" {
		fmt.Fprintf(w, "Hello, Welcome to the Main Page")
		return
	} else if !validFileType(w, fType) {
		return
	}

	fmt.Printf("File-path: %s\n", fullPath)
	fmt.Printf("File-name: %s\n", filename)
	fmt.Printf("File-type: %s\n", fType)

	if r.Method == "POST" {
		postHandler(w, r, fType)
	} else {
		getHandler(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------- <sending file> ----------")
	http.ServeFile(w, r, "./"+r.URL.String())
}

// Works except for corruption of photofiles
func postHandler(w http.ResponseWriter, r *http.Request, fType string) {
	fmt.Println("---------- <receiving file> ----------")

	f, error := os.Create(path.Base(r.URL.Path))
	if error != nil {
		fmt.Printf("Error on create file <POST>: %v\n", error)
		http.Error(w, "Error on create file <POST>: "+error.Error(), 500)

		return
	}

	switch fType {
	case "jpeg", "jpg", "gif":
		{
			img, _, error := image.Decode(r.Body)
			if error != nil {
				fmt.Printf("Error on decoding request data to image <POST>: %v\n", error)
				http.Error(w, "Error on decoding request data to image <POST>: "+error.Error(), 500)
				return
			}

			if fType == "gif" {
				error = gif.Encode(f, img, nil)
			} else {
				error = jpeg.Encode(f, img, nil)
			}

			if error != nil {
				fmt.Printf("Error on writing image to file <POST>: %v\n", error)
				http.Error(w, "Error on writing image to file <POST>: "+error.Error(), 500)
				return
			}

		}
	default:
		{
			_, error = io.Copy(f, r.Body)
			if error != nil {
				fmt.Printf("Error on writing to file <POST>: %v\n", error)
				r.Body.Close()
				return
			}
		}
	}

	r.Body.Close()
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

func validFileType(w http.ResponseWriter, fType string) bool {
	switch fType {
	case "html", "txt", "gif", "jpeg", "jpg", "css":
		return true
	default:
		http.Error(w, "File type not supported", 400)
		return false
	}
}

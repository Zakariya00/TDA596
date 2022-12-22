package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

/* HTTP Server Functions*/

// Gets rid of annoying favicon requests
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// serverHandler Main handler for http requests to server, with sync mechanism
func serverHandler(w http.ResponseWriter, r *http.Request) {

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

	if fileExists(path.Base(r.URL.Path)) {
		http.ServeFile(w, r, "./"+r.URL.String())
		return
	}
	proxyHandler(w, r)
}

// postHandler Handles POST requests
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

	/**/
	fileKey := hash(handler.Filename).String()
	belongsTo := cNode.find(cNode.LocalNode, fileKey)

	_, ok := cNode.Bucket[fileKey]
	_, ok1 := cNode.Backups[fileKey]
	if belongsTo.Id == cNode.LocalNode.Id || ok {
		fmt.Println("<Received file>")
		cNode.Bucket[fileKey] = handler.Filename
		return
	} else if ok1 {
		fmt.Println("Received file backup")
		cNode.Backups[fileKey] = handler.Filename
		return
	}

	fmt.Println("<Received then sent file>")
	postSender(belongsTo.Address, handler.Filename)
	/**/
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

// Serves get requests by forwarding them to the Main server and sends back response
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	forwardtoNode := cNode.lookUp(path.Base(r.URL.Path)).Address
	if forwardtoNode == cNode.LocalNode.Address {
		http.ServeFile(w, r, "./"+r.URL.String())
		return
	}

	forwardToServerURL := "http://" + forwardtoNode
	resp, err := http.Get(forwardToServerURL + r.URL.String())
	if err != nil {
		http.Error(w, "Failure in forwarding request: "+err.Error(), 400)
		return
	}

	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error in forwarding main server response: "+err.Error(), 500)
		return
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

// Copy source headers to destination header (use for sending back main server response)
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

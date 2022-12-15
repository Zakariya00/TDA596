package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// postSender sends files to nodes through http
func postSender(address string, filePath string) {
	uri := "http://" + address
	fileName := path.Base(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error on opening file:", err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("myFile", fileName)
	if err != nil {
		fmt.Println("Error on writing to form:", err)
		return
	}
	defer writer.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		fmt.Println("Failed to copy file for sending:", err)
		return
	}
	writer.Close()

	_, err = http.Post(uri, writer.FormDataContentType(), body) //
	if err != nil {
		fmt.Println("Failed to send Post:", err)
		return
	}
	/*
		res, err := httputil.DumpResponse(response, true)
		if err != nil {
			return
		}
		fmt.Print(string(res) + "\n")
	*/

	fmt.Println("<Sent file>")
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httputil"
	"os"
	"path"
)

var maxAttempts int = 5

// postSender sends files to nodes through http
func postSender(address string, filePath string) {
	for i := 0; i < maxAttempts; i++ {
		err := postSender1(address, filePath)
		if err != nil {
			continue
		}
		break
	}
}

func postSender1(address string, filePath string) error {
	client, err := httpsClient("secure_chord.crt")
	if err != nil {
		fmt.Println("Error on setting https client;", err)
		return err
	}
	newAddress, _ := setHttpsPort(address)
	uri := "https://" + newAddress
	fileName := path.Base(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error on opening file:", err)
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("myFile", fileName)
	if err != nil {
		fmt.Println("Error on writing to form:", err)
		return err
	}
	defer writer.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		fmt.Println("Failed to copy file for sending:", err)
		return err
	}
	writer.Close()

	//_, err = http.Post(uri, writer.FormDataContentType(), body) //
	response, err := client.Post(uri, writer.FormDataContentType(), body) //d
	if err != nil {
		fmt.Println("Failed to send Post:", err)
		return err
	}

	res, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}
	fmt.Print(string(res) + "\n")

	fmt.Println("<Sent file>")
	return nil
}

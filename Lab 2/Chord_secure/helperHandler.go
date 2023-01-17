package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

/* Helper Functions */

// getLocalAddress finds local ip address
func getLocalAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// fileExists checks if file exists locally
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	} else {
		fmt.Printf("File does not exist\n")
		return false
	}
}

// setHttpsPort increments address to the https port
func setHttpsPort(address string) (string, error) {
	u, err := url.Parse("http://" + address)
	if err != nil {
		fmt.Println("URL paring error:", err)
		return "", err
	}

	p, err := strconv.Atoi(u.Port())
	if err != nil {
		fmt.Println("URL paring error <port>:", err)
		return "", err
	}

	newAddress := u.Hostname() + ":" + strconv.Itoa(p+1) + u.Path
	return newAddress, nil
}

// httpsClient sets the https client
func httpsClient(crt string) (*http.Client, error) {
	f, err := os.Open(crt)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	caCert, err := os.ReadFile(crt)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},
	}
	return client, nil
}

// httpsServer sets the https server
func httpsServer(crt string) (*http.Server, error) {
	f, err := os.Open(crt)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	caCert, err := os.ReadFile(crt)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		//ClientAuth: tls.RequireAndVerifyClientCert,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  caCertPool,
	}
	srv := &http.Server{
		Addr:      "",
		Handler:   nil,
		TLSConfig: cfg,
	}
	return srv, nil
}

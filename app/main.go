package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		if conn != nil {
			fmt.Println("Connection established with client")
			go handleConn(conn)
		}
	}

}

func handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
			os.Exit(1)
		}
	}(conn)

	var request []byte
	request = make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	var requestData RequestData
	err = requestData.ParseRequest(request, n)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Request Data:")
	fmt.Println(requestData)

	response := getResponseBasedOnPath(requestData)
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}

	return
}

func getResponseBasedOnPath(requestData RequestData) []byte {
	switch requestData.Url {
	case "/":
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	case "/user-agent":
		userAgent := requestData.Headers["User-Agent"]
		contentLength := len(userAgent)
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, userAgent))
	default:
		return getResponseBasedOnFamilyPath(requestData)
	}
}

func getResponseBasedOnFamilyPath(requestData RequestData) []byte {
	pathSegments := strings.Split(requestData.Url, "/")

	for i, segment := range pathSegments {
		if strings.Contains(segment, "echo") {
			if i == len(pathSegments)-1 {
				return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
			}
			length := len(pathSegments[i+1])
			return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", length, pathSegments[i+1]))
		}
	}
	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}

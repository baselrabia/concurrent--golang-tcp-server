package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	handleConn(conn)
}
func handleConn(conn net.Conn) {
	reqStr := readConnToString(conn)
	fmt.Println(reqStr)
	httpReq := newHttpReq(reqStr)
	responder(conn, httpReq)
}
func readConnToString(conn net.Conn) string {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	return string(buff)
}
func responder(conn net.Conn, req *httpReq) {
	response := NotFoundResponse
	if req.path == "/" {
		response = OkResponse
	}
	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Println("Error writing into connection: ", err.Error())
		os.Exit(1)

	}
}

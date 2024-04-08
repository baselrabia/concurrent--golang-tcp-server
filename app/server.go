package main

import (
	"fmt"
	"strconv"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"github.com/codecrafters-io/http-server-starter-go/app/connection"

)

const PORT = 4221

var PORT_STR = strconv.Itoa(PORT)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:"+PORT_STR)
	if err != nil {
		fmt.Println("Failed to bind to port " + PORT_STR)
		os.Exit(1)
	}
	
	// connChan := make(chan net.Conn)

	//	can handle concurrent requests
	// go connection.HandleConnections(connChan)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		// connChan <- conn
		go connection.HandleConn(conn)
	}
}

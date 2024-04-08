package connection

import (
	"fmt"
	"net"
	"time"
	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

func HandleConnections(connChan chan net.Conn) {
	for {
		conn := <-connChan
		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error while closing connection: ", err.Error())
		}
	}()
	_ = conn.SetDeadline(time.Now().Add(15 * time.Second))
	reqBuf := make([]byte, 1024)
	bytesRead, err := conn.Read(reqBuf)
	if err != nil {
		fmt.Println("Error while reading from connection: ", err.Error())
		conn.Close()
		return
	}
	reqBuf = reqBuf[:bytesRead]
	
	fmt.Printf("Bytes read: %d\n", bytesRead)
	fmt.Printf("Request:\n---Request Start---\n%s\n---Request End---\n", reqBuf)

	fmt.Println("parsing request...")
	req, _ := request.ParseRequest(reqBuf)
	fmt.Println("creating response...")
	res := response.FromRequest(req)

	str, err := res.String()
	if err != nil {
		fmt.Println("Error converting response to string: ", err.Error())
		return
	}
	bytes := []byte(str)
	bytesWritten, err := conn.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	if bytesWritten != len(bytes) {
		fmt.Printf("Did not write proper amount of bytes. Expected: %d, Got: %d\n", len(bytes), bytesWritten)
	}
	fmt.Printf("Wrote %d bytes\n", bytesWritten)
	if bytesWritten > 0 {
		fmt.Printf("Response:\n---RESPONSE Start---\n%s\n---RESPONSE End---\n", bytes[:bytesWritten])
 
	}
}
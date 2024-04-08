package response

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
)

type Response struct {
	Status  int
	Headers map[string]string

	Body []byte
	str  strings.Builder
}

const HTTP_VER = "HTTP/1.1"
const LINE_BREAK = "\r\n"

const (
	OkResponse       = "HTTP/1.1 200 OK\r\n\r\n"
	NotFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
)

var fileServerDir = func() string {
	idx := slices.Index(os.Args, "--directory")
	if idx == -1 || idx == len(os.Args)-1 {
		fmt.Println("WARNING: No valid '--directory' provided. Cannot serve files")
		return ""
	}
	return os.Args[idx+1]

}()

func FromRequest(req *request.Request) *Response {
	// In descending priority.
	const ECHO_PRE = "/echo/"
	const FILE_PRE = "/files/"
	const HEADER_PRE = "/"

	switch {
	case req == nil:
		return &Response{Status: 404}
	case req.Path == "/":
		return &Response{Status: 200}

	case strings.HasPrefix(req.Path, ECHO_PRE):
		body := req.Path[len(ECHO_PRE):]
		res := &Response{Status: 200}
		res.addBody("text/plain", []byte(body))
		return res

	// Serve file contents route
	case strings.HasPrefix(req.Path, FILE_PRE):
		filename := req.Path[len(FILE_PRE):]
		entries, err := os.ReadDir(fileServerDir)
		if err != nil {
			fmt.Printf("Could read file server directory %q. Cannot serve file %q: %s\n", fileServerDir, filename, err.Error())
			return &Response{Status: 404}
		}
		path := filepath.Join(fileServerDir, filename)
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() != filename {
				continue
			}
			contents, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Could not read file %q: %s\n", path, err.Error())
				return &Response{Status: 400}
			}
			res := &Response{Status: 200}
			res.addBody("application/octet-stream", contents)
		  
			return res
		}
		fmt.Printf("Could not find file. Filename: %q. Path: %q\n", filename, path)
	// Echo header route
	case strings.HasPrefix(req.Path, HEADER_PRE):
		header := strings.ToLower(req.Path[len(HEADER_PRE):])

		for key, value := range req.Headers {
			if strings.ToLower(key) != header {
				continue
			}
			res := &Response{Status: 200}
			res.addBody("text/plain", []byte(value))
			return res
		}
	}

	return &Response{Status: 404}
}

func (res *Response) addBody(contentType string, body []byte) {
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["Content-Type"] = contentType
	res.Headers["Content-Length"] = strconv.Itoa(len(body))

	res.Body = body
}

func (res *Response) String() (string, error) {
	res.writeStatus()
	res.writeHeaders()

	res.writeBody()
	return res.str.String(), nil
}

func (res *Response) writeStatus() {
	res.str.WriteString(HTTP_VER)
	res.str.WriteRune(' ')
	res.str.WriteString(strconv.Itoa(res.Status))
	status := ""
	switch res.Status {
	case 200:
		status = "OK"
	case 404:
		status = "Not Found"
	}
	if status != "" {
		res.str.WriteRune(' ')
		res.str.WriteString(status)
	}
	res.str.WriteString(LINE_BREAK)
}

func (res *Response) writeHeaders() {

	for key, val := range res.Headers {
		res.str.WriteString(key)
		res.str.WriteString(": ")
		res.str.WriteString(val)
		res.str.WriteString(LINE_BREAK)
	}
	res.str.WriteString(LINE_BREAK)
}

func (res *Response) writeBody() {
	res.str.Write(res.Body)
	// TODO: Might have to add \r\n here?
}

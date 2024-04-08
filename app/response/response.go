package response

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"

)

type Response struct {
	Status  int
	Headers map[string]string

	Body string
	str  strings.Builder
}

const HTTP_VER = "HTTP/1.1"
const LINE_BREAK = "\r\n"

const (
	OkResponse       = "HTTP/1.1 200 OK\r\n\r\n"
	NotFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
)

func FromRequest(req *request.Request) *Response {
	const ECHO_PRE = "/echo/"
	const HEADER_PRE = "/"

	switch {
	case req == nil:
		return &Response{Status: 404}
	case req.Path == "/":
		return &Response{Status: 200}
	case strings.HasPrefix(req.Path, ECHO_PRE):
		body := req.Path[len(ECHO_PRE):]
		res := &Response{Status: 200}
		res.addBody("text/plain", body)
		return res
	case strings.HasPrefix(req.Path, HEADER_PRE):
		header := strings.ToLower(req.Path[len(HEADER_PRE):])
		
		for key, value := range req.Headers {
			if strings.ToLower(key) != header {
				continue
			}
			res := &Response{Status: 200}
			res.addBody("text/plain", value)
			return res
		}
	}

	return &Response{Status: 404}
}

func (res *Response) addBody(contentType string, body string) {
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
	res.str.WriteString(res.Body)
	// TODO: Might have to add \r\n here?
}

package request

import (
	"fmt"
	"strings"
	"unicode"
)

type Parser struct {
	input string
	pos   int
	req   *Request
}
type Request struct {
	Method  string
	Path    string
	HttpVer string
	Headers map[string]string
	Body    []byte
}

const LINE_BREAK = "\r\n"

func ParseRequest(input []byte) (*Request, error) {
	str := string(input)
	req := &Request{}
	parser := Parser{input: str[:], pos: 0, req: req}
	fmt.Println("parsing verb...")
	if ok := parser.parseVerb(); !ok {
		return nil, nil
	}
	fmt.Println("parsing headers...")
	parser.parseHeaders()
	fmt.Println("parsing body...")
	parser.parseBody()

	return req, nil
}

func (p *Parser) parseVerb() bool {
	lineBreakIdx := strings.Index(p.input[p.pos:], LINE_BREAK)
	if lineBreakIdx == -1 {
		return false
	}
	line := p.input[p.pos:lineBreakIdx]
	tokens := strings.Split(line, " ")
	p.req.Method = tokens[0]
	p.req.Path = tokens[1]
	p.req.HttpVer = tokens[2]
	p.pos = len(line) + len(LINE_BREAK)

	if p.pos >= len(p.input) {
		p.pos = len(p.input) - 1
	}

	return true
}
func (p *Parser) parseHeaders() {
	if p.req.Headers == nil {
		p.req.Headers = make(map[string]string)
	}
	
	fmt.Printf("Req so far: %q\n", p.input[p.pos:])
	endMarker := LINE_BREAK + LINE_BREAK
	endIdx := strings.Index(p.input[p.pos:], endMarker)
	fmt.Printf("endIdx: %d. endMarker: %q\n", endIdx, endMarker)
	if endIdx == -1 {
		fmt.Println("No end to headers in request!")
		panic("No end to headers in request!")
	}
	afterEndIdx := p.pos + endIdx + len(endMarker)
	fmt.Println("b")
	headers := p.input[p.pos:afterEndIdx]
	fmt.Println("c")
	p.pos = afterEndIdx // Not checking for overflows
	fmt.Println("a")
	lines := strings.Split(headers, LINE_BREAK)


	for _, line := range lines {
		if line == "" {
			continue
		}
		str := strings.TrimLeftFunc(line, unicode.IsSpace)
		sepIdx := strings.IndexRune(str, ':')
		if sepIdx == -1 || sepIdx+1 >= len(str) {
			continue
		}
		key := str[:sepIdx]
		value := strings.TrimSpace(str[sepIdx+1:])
		p.req.Headers[key] = value
	}
}

func (p *Parser) parseBody() {
	if p.pos >= len(p.input) {
		fmt.Println("No body provided in request")
		return
	}
	p.req.Body = []byte(p.input[p.pos:])
	fmt.Println("Body:", string(p.req.Body), "Len:", len(p.req.Body))
}

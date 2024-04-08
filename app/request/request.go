package request
import (
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
}
 
const LINE_BREAK = "\r\n"
func ParseRequest(input []byte) (*Request, error) {
	str := string(input)
	req := &Request{}
	parser := Parser{input: str[:], pos: 0, req: req}
	if ok := parser.parseVerb(); !ok {
		return nil, nil
	}
	parser.parseHeaders()
 
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
	if p.pos > len(p.input) {
		p.pos = len(p.input) - 1
	}
 	return true
}
func (p *Parser) parseHeaders() {
	if p.req.Headers == nil {
		p.req.Headers = make(map[string]string)
	}
	lines := strings.Split(p.input[p.pos:], LINE_BREAK)
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
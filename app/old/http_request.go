package old

import (
	"bufio"
	"fmt"
	"strings"
)

type httpReq struct {
	method      string
	path        string
	httpVersion string
	headers map[string]string
}

func newHttpReq(reqStr string) *httpReq {
	return parseReq(reqStr)
}

func parseReq(reqStr string) *httpReq {
	httpReq := new(httpReq)
	reqScanner := bufio.NewScanner(strings.NewReader(reqStr))
	reqScanner.Split(bufio.ScanLines)
	reqScanner.Scan()
	startLine := reqScanner.Text()
	httpReq.parseStartLine(startLine)
	return httpReq
}

func (h *httpReq) parseStartLine(startLine string) {
	fmt.Println("sssssss" ,startLine)
	s := bufio.NewScanner(strings.NewReader(startLine))
	s.Split(bufio.ScanWords)
	s.Scan()
	h.method = s.Text()
	s.Scan()
	h.path = s.Text()
	// s.Scan()
	// h.headers = s.Text()

	h.httpVersion = s.Text()
}

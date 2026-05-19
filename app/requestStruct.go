package main

import (
	"fmt"
	"strings"
)

var CRLF = "\r\n"

type RequestData struct {
	Method   string
	Url      string
	Protocol string
	Lines    string
	Headers  map[string]string
	Body     string
}

func (r *RequestData) ParseRequest(request []byte, requestLength int) error {
	buffer := request[:requestLength]
	lines := strings.Split(string(buffer), CRLF)

	parsingHeader := true
	r.Headers = make(map[string]string)
	for index, line := range lines {
		if index == 0 { // request line
			err := r.ParseLines(line)
			if err != nil {
				return err
			}
		} else if parsingHeader { //header
			if line == "" {
				parsingHeader = false
			} else {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid header format at line %d", index)
				}
				for i, part := range parts {
					parts[i] = strings.TrimSpace(part)
				}
				r.Headers[parts[0]] = parts[1]
			}
		} else { // body
			r.Body += line + CRLF
		}
	}
	return nil
}

func (r *RequestData) ParseLines(lines string) error {
	if lines == "" {
		fmt.Println("Invalid request format")
		return fmt.Errorf("invalid request format")
	}

	parts := strings.Fields(lines)
	if len(parts) != 3 {
		fmt.Println("Invalid request format")
		return fmt.Errorf("invalid request format")
	}
	r.Method = parts[0]
	r.Url = parts[1]
	r.Protocol = parts[2]
	return nil
}

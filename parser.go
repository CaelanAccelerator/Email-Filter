package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Email struct {
	Headers map[string]string
	Body    string
}

func parseEmailStream(r io.Reader) Email {
	email := Email{Headers: make(map[string]string)}
	scanner := bufio.NewScanner(r)
	var bodyLines []string
	isHead := true

	for scanner.Scan() {
		line := trimCR(scanner.Text())

		if isHead {
			// if is empty line, following content is body
			if line == "" {
				isHead = false
				continue
			}

			hKey, hVal, hOK := splitHeader(line)
			if hOK {
				email.Headers[hKey] = hVal
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Scan Error: ", scanner.Err())
	}

	email.Body = strings.Join(bodyLines, "\n")
	return email
}

func trimCR(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != '\r' {
			return s[:i+1]
		}
	}
	return ""
}

func splitHeader(line string) (key string, value string, ok bool) {
	idx := strings.Index(line, ":")
	if idx == -1 {
		return "", "", false
	}

	key = strings.ToLower(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	ok = true
	return
}

package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer ln.Close()
	fmt.Println("listening on :8080")

	serve(ln)
}

func serve(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	email := parseEmailStream(conn)
	verdict := classify(email)

	record := struct {
		Email   Email
		Verdict Verdict
	}{
		Email:   email,
		Verdict: verdict,
	}
	data, _ := json.MarshalIndent(record, "", "  ")
	fmt.Println(string(data))

	fmt.Fprintf(conn, "from: %s\n subject: %s\n is recieved", email.Headers["from"], email.Headers["subject"])

	err := persistVerdict(email, verdict)
	if err != nil {
		fmt.Println("Persist error: ", err)
	}
}

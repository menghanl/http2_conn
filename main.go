package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"playground/http2_conn/httpconn"

	log "github.com/sirupsen/logrus"
)

func server() {
	listener := httpconn.Listen(":4430")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go io.Copy(conn, conn)
	}
}

func main() {
	go server()
	go client()
	go client()
	select {}
}

func client() {
	dialer := &httpconn.Dialer{
		InsecureSkipVerify: true,
	}
	conn := dialer.Dial("localhost:4430")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(conn, "It is now %v\n", time.Now())
		}
	}()

	go func() {
		n, err := io.Copy(os.Stdout, conn)
		log.Fatalf("copied %d, %v", n, err)
	}()
}

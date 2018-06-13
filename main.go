package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"playground/http2_conn/httpconn"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func server() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "MAGIC" {
			http.Error(w, "MAGIC required. Not "+r.Method, 400)
			return
		}
		io.Copy(flushWriter{w}, r.Body)
	}

	var srv http.Server
	srv.Addr = ":4430"
	http.HandleFunc("/ECHO", handler)
	http2.ConfigureServer(&srv, &http2.Server{})
	log.Fatal(srv.ListenAndServeTLS("./tls/server1.pem", "./tls/server1.key"))
	// log.Fatal(srv.ListenAndServe())
}

func main() {
	go server()
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

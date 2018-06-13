package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
		fmt.Println("flushing")
	}
	return
}

func server() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(flushWriter{w}, r.Body)
	}

	http.HandleFunc("/fakeconn", handler)
	var srv http.Server
	srv.Addr = ":8080"
	http2.ConfigureServer(&srv, &http2.Server{})
	srv.ListenAndServeTLS("./tls/server1.pem", "./tls/server1.key")
}

func main() {
	go server()
	go client()
	select {}
}

func client() {
	t := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: t,
	}
	pr, pw := io.Pipe()
	req, err := http.NewRequest("PUT", "https://localhost:8080/fakeconn", ioutil.NopCloser(pr))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(pw, "It is now %v\n", time.Now())
		}
	}()
	go func() {
		res, err := c.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Got: %#v", res)
		n, err := io.Copy(os.Stdout, res.Body)
		log.Fatalf("copied %d, %v", n, err)
	}()
}

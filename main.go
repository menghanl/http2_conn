package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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
		if r.Method != "PUT" {
			http.Error(w, "PUT required.", 400)
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
	t := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: t,
	}
	pr, pw := io.Pipe()
	req, err := http.NewRequest("PUT", "https://localhost:4430/ECHO", ioutil.NopCloser(pr))
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

package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

type capitalizeReader struct {
	r io.Reader
}

func (cr capitalizeReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	for i, b := range p[:n] {
		if b >= 'a' && b <= 'z' {
			p[i] = b - ('a' - 'A')
		}
	}
	return
}

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

func echoCapitalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "PUT required.", 400)
		return
	}
	io.Copy(flushWriter{w}, capitalizeReader{r.Body})
}

func registerHandlers() {
	mux2 := http.NewServeMux()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mux2.ServeHTTP(w, r)
	})
	mux2.HandleFunc("/ECHO", echoCapitalHandler)
}

func main() {
	var srv http.Server
	srv.Addr = ":4430"

	registerHandlers()

	http2.ConfigureServer(&srv, &http2.Server{})

	go func() {
		log.Fatal(srv.ListenAndServeTLS("./tls/server1.pem", "./tls/server1.key"))
	}()
	select {}
}

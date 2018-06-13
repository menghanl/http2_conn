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
		// go func() {
		// for {
		// 	w.Write([]byte("hello aaa"))
		// 	w.(http.Flusher).Flush()
		// 	time.Sleep(time.Second)
		// }
		// }()
		// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
		io.Copy(os.Stdout, r.Body)
	}

	http.HandleFunc("/fakeconn", handler)
	// log.Fatal(http.ListenAndServeTLS(":8080", "./tls/server1.pem", "./tls/server1.key", nil))

	var srv http.Server
	srv.Addr = ":8080"
	http2.ConfigureServer(&srv, &http2.Server{})
	srv.ListenAndServeTLS("./tls/server1.pem", "./tls/server1.key")
}

func client() {

	t := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: t,
	}
	pr, pw := io.Pipe()
	req, err := http.NewRequest("PUT", "https://localhost:8080/ECHO", ioutil.NopCloser(pr))
	go func() {
		for {
			fmt.Fprintf(pw, "It is now %v\n", time.Now())
			time.Sleep(time.Second)
		}
	}()
	// req, err = http.NewRequest("PUT", "https://localhost:8080/fakeconn", bytes.NewBufferString("aaa"))
	res, err := c.Do(req)
	fmt.Println("after Do")
	fmt.Println(res, err)
	io.Copy(os.Stdout, res.Body)

}

func main() {
	// fmt.Println("hello")

	// // go server()
	// go client()

	// select {}

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
	select {}
}

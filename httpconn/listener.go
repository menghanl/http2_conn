package httpconn

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
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
	}
	return
}

// Listen creates a listener that creates conn on top of http2.
func Listen(addr string) net.Listener {
	lis := &listener{
		connCh: make(chan net.Conn),
	}

	var srv http.Server
	srv.Addr = addr // ":4430"
	http.HandleFunc("/ECHO", lis.httpHandler)
	http2.ConfigureServer(&srv, &http2.Server{})
	go func() {
		log.Fatal(srv.ListenAndServeTLS("./tls/server1.pem", "./tls/server1.key"))
	}()

	return lis
}

// listener implements net.Listener, and creates conn that wrapps around http2.
type listener struct {
	closeOnce sync.Once
	connCh    chan net.Conn
}

func (l *listener) httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler called")
	if r.Method != "MAGIC" {
		http.Error(w, "MAGIC required. Not "+r.Method, 400)
		return
	}
	conn := &serverConn{
		out: flushWriter{w},
		in:  r.Body,
	}
	fmt.Println("writing to channel")
	l.connCh <- conn
	fmt.Println("wrote to channel")
	select {}
}

func (l *listener) Accept() (net.Conn, error) {
	conn, ok := <-l.connCh
	if !ok {
		return nil, fmt.Errorf("listener closed")
	}
	//
	// The code blocks indefinitely if the following read() tries to read the
	// whole magicHandshakeStr.
	//
	magicBytes := make([]byte, len(magicHandshakeStr)-1)
	// if _, err := io.ReadFull(conn, magicBytes); err != nil {
	if _, err := conn.Read(magicBytes); err != nil {
		return nil, fmt.Errorf("failed to handshake: %v", err)
	}
	return conn, nil
}

func (l *listener) Close() error {
	l.closeOnce.Do(func() {
		close(l.connCh)
	})
	return nil
}

func (l *listener) Addr() net.Addr {
	return constFakeAddr
}

// serverConn implements net.Conn.
type serverConn struct {
	out io.Writer
	in  io.ReadCloser
}

func (c *serverConn) Read(b []byte) (n int, err error) {
	return c.in.Read(b)
}

func (c *serverConn) Write(b []byte) (n int, err error) {
	return c.out.Write(b)
}

func (c *serverConn) Close() error {
	return c.in.Close()
}

func (c *serverConn) LocalAddr() net.Addr                { return constFakeAddr }
func (c *serverConn) RemoteAddr() net.Addr               { return constFakeAddr }
func (c *serverConn) SetDeadline(t time.Time) error      { return nil }
func (c *serverConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *serverConn) SetWriteDeadline(t time.Time) error { return nil }

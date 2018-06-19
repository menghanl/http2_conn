package httpconn

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// Dialer implements net.Dialer, and creates conn that wrapps around http2.
type Dialer struct {
	InsecureSkipVerify bool
}

// Dial creates a conn wrapper on top of a http2 client.
func (d *Dialer) Dial(target string) net.Conn {
	c := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: d.InsecureSkipVerify},
		},
	}
	pr, pw := io.Pipe()
	req, err := http.NewRequest("MAGIC", "https://"+target+"/ECHO", ioutil.NopCloser(pr))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// Need to send something as a handshake, otherwise the following Do()
		// will block. Don't know why...
		//
		// Also, if the server tries to read the whole string, it still blocks
		// (maybe also in Do()). Currently server reads one byte less, but that
		// needs to be fixed.
		pw.Write([]byte(magicHandshakeStr))
	}()
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got: %#v", res)

	return &clientConn{
		out: pw,
		in:  res.Body,
	}
}

// clientConn implements net.Conn.
type clientConn struct {
	out io.Writer
	in  io.ReadCloser
}

func (c *clientConn) Read(b []byte) (n int, err error) {
	return c.in.Read(b)
}

func (c *clientConn) Write(b []byte) (n int, err error) {
	return c.out.Write(b)
}

func (c *clientConn) Close() error {
	return c.in.Close()
}

func (c *clientConn) LocalAddr() net.Addr                { return constFakeAddr }
func (c *clientConn) RemoteAddr() net.Addr               { return constFakeAddr }
func (c *clientConn) SetDeadline(t time.Time) error      { return nil }
func (c *clientConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *clientConn) SetWriteDeadline(t time.Time) error { return nil }

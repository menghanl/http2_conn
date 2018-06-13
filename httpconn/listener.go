package httpconn

import (
	"fmt"
	"net"
	"time"
)

// Listener implements net.Listener, and creates conn that wrapps around http2.
type Listener struct{}

// Accept returns a conn wrapper on top of http2 stream.
func (l *Listener) Accept() net.Conn {
	return &serverConn{}
}

// serverConn implements net.Conn.
type serverConn struct{}

func (c *serverConn) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *serverConn) Write(b []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *serverConn) Close() error {
	return fmt.Errorf("not implemented")
}

func (c *serverConn) LocalAddr() net.Addr                { return constFakeAddr }
func (c *serverConn) RemoteAddr() net.Addr               { return constFakeAddr }
func (c *serverConn) SetDeadline(t time.Time) error      { return nil }
func (c *serverConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *serverConn) SetWriteDeadline(t time.Time) error { return nil }

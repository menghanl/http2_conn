package httpconn

import (
	"fmt"
	"net"
	"time"
)

// Dialer implements net.Dialer, and creates conn that wrapps around http2.
type Dialer struct{}

// Dial creates a conn wrapper on top of a http2 client.
func (d *Dialer) Dial(target string) net.Conn {
	return &clientConn{}
}

// clientConn implements net.Conn.
type clientConn struct{}

func (c *clientConn) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *clientConn) Write(b []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *clientConn) Close() error {
	return fmt.Errorf("not implemented")
}

func (c *clientConn) LocalAddr() net.Addr                { return constFakeAddr }
func (c *clientConn) RemoteAddr() net.Addr               { return constFakeAddr }
func (c *clientConn) SetDeadline(t time.Time) error      { return nil }
func (c *clientConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *clientConn) SetWriteDeadline(t time.Time) error { return nil }

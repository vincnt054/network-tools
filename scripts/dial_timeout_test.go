package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func DialTimeOut(network, address string, timeout time.Duration,
) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_,  addr string, _ syscall.RawConn) error { // Overrides Control function to return error
			return &net.DNSError {
				Err: "connection timed out",
					Name: addr,
					Server: "127.0.0.1",
					IsTimeout: true,
					IsTemporary: true,
				}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeOut(t *testing.T) {
	c, err := DialTimeOut("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}

	if !nErr.Timeout() {
		t.Fatal("error is not a timeout")
	}
}

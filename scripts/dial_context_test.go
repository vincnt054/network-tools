package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(5 * time.Second) // Create context with deadline of 5 seconds
	ctx, cancel := context.WithDeadline(context.Background(), dl) // Create the context and its cancel function
	defer cancel() // elligible for garbage collected

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// Sleep to reach context's deadline.
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	} // Override dialer's Control function

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}
	if ctx.Err() != context.DeadlineExceeded { // sanity check at the end of the test
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}

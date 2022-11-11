package ch03

import (
	"context"
	"fmt"
	"io"
	"time"
)

func ExamplePinger() {
	ctx, cancel := context.WithCancel(context.Background())
	r, w = io.Pipe()
	done := make(chan struct{})
	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receive

package echo

import (
	"context"
	"net"
)

func echoServerUDP(ctrx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)

		for {
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}

func streamingEchoServer(ctx context.Context, network string, addr string) (net.Addr, error) {
	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for {
			conn, err := s.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() { _ = conn.Close() }()

				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	return s.Addr(), nil
}

func datagramEchoServer(ctx context.Context, network string, addr string) (net.Addr, error) {
	s, err := net.ListenPacket(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<- ctx.Done()
			_ = s.Close()
			if network == "unixgram" {
				_ = os.Remove(addr)
			}
		}()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}
		}
	}()
	return s.LocalAddr(), nil
}

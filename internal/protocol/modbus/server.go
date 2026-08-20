package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Handler interface {
	Handle(context.Context, Frame) (Frame, error)
}
type Server struct {
	Addr         string
	Handler      Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (s Server) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	go func() { <-ctx.Done(); listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return err
			}
		}
		go s.handle(ctx, conn)
	}
}
func (s Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		_ = conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		header := make([]byte, 7)
		if _, err := ioReadFull(conn, header); err != nil {
			return
		}
		length := int(binary.BigEndian.Uint16(header[4:6]))
		if length < 2 || length > 253 {
			return
		}
		body := make([]byte, length-1)
		if _, err := ioReadFull(conn, body); err != nil {
			return
		}
		raw := append(header, body...)
		frame, err := Decode(raw)
		if err != nil {
			return
		}
		reply, err := s.Handler.Handle(ctx, frame)
		if err != nil {
			return
		}
		encoded, err := Encode(reply)
		if err != nil {
			return
		}
		_ = conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
		if _, err = conn.Write(encoded); err != nil {
			return
		}
	}
}
func ioReadFull(c net.Conn, b []byte) (int, error) {
	total := 0
	for total < len(b) {
		n, err := c.Read(b[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

type EchoHandler struct{}

func (EchoHandler) Handle(_ context.Context, f Frame) (Frame, error) {
	if f.Function&0x80 != 0 {
		return f, fmt.Errorf("exception request")
	}
	return Frame{Transaction: f.Transaction, Unit: f.Unit, Function: f.Function, Payload: f.Payload}, nil
}

package transports

import (
	"fmt"
	"log"
	"net"

	"github.com/jtarchie/syslog/pkg/log"
)

type Writer interface {
	Write(*syslog.Log) error
}

type UDPServer struct {
	writer   Writer
	listener net.PacketConn
}

func NewUDPServer(port int, w Writer) (*UDPServer, error) {
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("could not start server: %s", err)
	}
	return &UDPServer{
		writer:   w,
		listener: listener,
	}, nil
}

func (s *UDPServer) Start() error {
	log.Printf("udp: starting server on addr %s", s.listener.LocalAddr().String())
	defer s.listener.Close()

	buffer := make([]byte, 1024)
	for {
		n, _, err := s.listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("could not read from UDP: %s", err)
			continue
		}

		parsed, err := syslog.Parse(buffer[:n])
		if err != nil {
			log.Printf("could not parse msg: %s", err)
			continue
		}
		s.writer.Write(parsed)
	}
	return nil
}

func (s *UDPServer) Addr() net.Addr {
	return s.listener.LocalAddr()
}

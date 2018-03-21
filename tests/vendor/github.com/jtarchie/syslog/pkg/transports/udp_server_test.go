package transports_test

import (
	"log"
	"net"
	"sync"
	"testing"

	"github.com/jtarchie/syslog/pkg/log"
	"github.com/jtarchie/syslog/pkg/transports"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func SendMessage(s *transports.UDPServer) {
	payload := []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)

	conn, err := net.Dial("udp", s.Addr().String())
	Expect(err).ToNot(HaveOccurred())
	defer conn.Close()

	_, err = conn.Write(payload)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Server", func() {
	It("accepts datagrams via UDP", func() {
		writer := &SpyWriter{}
		server, err := transports.NewUDPServer(0, writer)
		Expect(err).ToNot(HaveOccurred())
		go server.Start()

		SendMessage(server)

		Eventually(func() int {
			return len(writer.Logs())
		}).Should(Equal(1))

		Expect(writer.Logs()[0].Message()).To(BeEquivalentTo(`'su root' failed for lonvick on /dev/pts/8`))
	})
})

type SpyWriter struct {
	logs []*syslog.Log
	mu   sync.Mutex
}

func (s *SpyWriter) Write(log *syslog.Log) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = append(s.logs, log)
	return nil
}

func (s *SpyWriter) Logs() []*syslog.Log {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.logs
}

func TestTransports(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transports Suite")
}

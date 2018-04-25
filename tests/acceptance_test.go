package windows_syslog_acceptance_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Forwarding Loglines", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/tcp.yml")
	})

	AfterEach(func() {
		Cleanup()
	})

	It("forwards logs from /var/vcap/sys/log ", func() {
		message := counterString(500, "A")
		Eventually(WriteToTestFile(message)).Should(ContainSubstring(message))
	})
})

var _ = Describe("Forwarding Loglines using tls", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/tls.yml")
	})

	AfterEach(func() {
		Cleanup()
	})

	It("forwards logs from /var/vcap/sys/log ", func() {
		message := counterString(500, "A")
		Eventually(WriteToTestFile(message)).Should(ContainSubstring(message))
	})
})

func counterString(l int, s string) string {
	counterstring := ""
	for len(counterstring) < l {
		counterstring += s
		counterstring += strconv.Itoa(len(counterstring))
	}

	return counterstring[:l]
}

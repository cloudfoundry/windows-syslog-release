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

var _ = Describe("Disabling the log forwarder", func() {
	Context("when a storer is available and configuration is provided", func() {
		BeforeEach(func() {
			Cleanup()
			Deploy("manifests/disabled.yml")
		})
		AfterEach(func() {
			Cleanup()
		})
		It("doesn't forward logs from /var/vcap/sys/log ", func() {
			message := counterString(500, "A")
			Consistently(WriteToTestFile(message)).ShouldNot(ContainSubstring(message))
		})
	})

	Context("when no configuration is provided", func() {
		BeforeEach(func() {
			Cleanup()
		})
		AfterEach(func() {
			Cleanup()
		})

		It("Deploys disabled", func() {
			Deploy("manifests/disabled-no-config.yml")
		})
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

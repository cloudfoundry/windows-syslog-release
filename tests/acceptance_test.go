package windows_syslog_acceptance_test

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Forwarding Loglines", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/tcp.yml")
	})

	It("forwards logs from /var/vcap/sys/log ", func() {
		message := counterString(500, "A")
		Eventually(WriteToTestFile(message)).Should(ContainSubstring(message))
	})

	It("annotates logs with structured data specific to the instance", func() {
		ExpectedStructuredDataRegexp := fmt.Sprintf("\\[instance@47450 az=\".*\" deployment=\"%s\" director=\"test-env\" group=\"forwarder\" id=\".*\"\\]", DeploymentName())
		message := counterString(500, "A")
		Eventually(WriteToTestFile(message)).Should(MatchRegexp(ExpectedStructuredDataRegexp))
	})
})

var _ = Describe("Forwards windows event logs", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/events.yml")
	})

	It("forwards windows event logs", func() {
		Eventually(SSHForAccessLog()).Should(ContainSubstring("Microsoft-Windows-Security-Auditing"))
	})
})

var _ = Describe("Forwards windows event logs", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/events-custom.yml")
	})

	It("forwards windows event logs", func() {
		message := "AAAAAAAAAAAAAAAAAAAAAAAA"
		RegisterEventLog()
		Eventually(WriteMachineEventLog(message)).Should(ContainSubstring(message))
	})
})

var _ = Describe("Forwarding Loglines using tls", func() {
	BeforeEach(func() {
		Cleanup()
		Deploy("manifests/tls.yml")
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
		It("doesn't forward logs from /var/vcap/sys/log ", func() {
			message := counterString(500, "A")
			Consistently(WriteToTestFile(message)).ShouldNot(ContainSubstring(message))
		})
	})

	Context("when no configuration is provided", func() {
		BeforeEach(func() {
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

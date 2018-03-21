package syslog_test

import (
	"time"

	"github.com/jtarchie/syslog/pkg/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log Parse", func() {
	var (
		payload []byte
		log     *syslog.Log
	)

	BeforeEach(func() {
		payload = []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)
		log = &syslog.Log{}
	})

	It("parses valid messages", func() {
		_, err := syslog.Parse(payload)
		Expect(err).ToNot(HaveOccurred())
	})

	It("sets the version", func() {
		log, _ := syslog.Parse(payload)
		Expect(log.Version()).To(Equal(1))
	})

	Context("with the priority", func() {
		It("returns the severity", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Severity()).To(Equal(2))
		})

		It("returns the facility", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Facility()).To(Equal(4))
		})

		It("sets the priority", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Priority()).To(Equal(34))
		})
	})

	Context("with the timestamp", func() {
		It("returns a valid date object", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Timestamp().String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
		})

		Context("with the example timestamps from the RFC", func() {
			It("can parse them all", func() {
				timestamp := func(t string) time.Time {
					payload := []byte("<34>1 " + t + " mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
					log, err := syslog.Parse(payload)
					Expect(err).ToNot(HaveOccurred())

					return log.Timestamp()
				}

				Expect(timestamp("2003-10-11T22:14:15.00003Z").String()).To(Equal("2003-10-11 22:14:15.00003 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52Z").String()).To(Equal("1985-04-12 23:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52+00:00").String()).To(Equal("1985-04-12 23:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52+02:00").String()).To(Equal("1985-04-12 21:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T18:20:50.52-02:00").String()).To(Equal("1985-04-12 20:20:50.52 +0000 UTC"))
			})

			It("fails parsing on unsupported formats", func() {
				payload := []byte("<34>1 2003-10-11T22:14:15.003Z07:00 mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
				_, err := syslog.Parse(payload)
				Expect(err).To(HaveOccurred())
			})
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.Timestamp().IsZero()).To(BeTrue())
		})
	})

	Context("with the hostname", func() {
		It("sets the hostname", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Hostname()).To(BeEquivalentTo("mymachine.example.com"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.Hostname()).To(BeNil())
		})
	})

	Context("with the app name", func() {
		It("sets the app name", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Appname()).To(BeEquivalentTo("su"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.Appname()).To(BeNil())
		})
	})

	Context("with the proc id", func() {
		It("sets the proc id", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.ProcID()).To(BeEquivalentTo("12345"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.ProcID()).To(BeNil())
		})
	})

	Context("with the log id", func() {
		It("sets the log id", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.MsgID()).To(BeEquivalentTo("98765"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.MsgID()).To(BeNil())
		})
	})

	Context("with structure data", func() {
		It("sets structure data", func() {
			log, _ := syslog.Parse(payload)
			data := *log.StructureData()
			Expect(data.ID()).To(BeEquivalentTo("exampleSDID@32473"))
			Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
				{[]byte("iut"), []byte("3")},
				{[]byte("eventSource"), []byte("Application")},
				{[]byte("eventID"), []byte("1011")},
			}))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _ := syslog.Parse(payload)
			Expect(log.StructureData()).To(BeNil())
		})
	})

	Context("with a message", func() {
		It("sets the message", func() {
			log, _ := syslog.Parse(payload)
			Expect(log.Message()).To(BeEquivalentTo("'su root' failed for lonvick on /dev/pts/8"))
		})

		It("sets nil for no message", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - -")
			log, _ := syslog.Parse(payload)
			Expect(log.Message()).To(BeNil())
		})
	})
})

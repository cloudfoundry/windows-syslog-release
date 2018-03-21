package syslog_test

import (
	"testing"

	"github.com/jtarchie/syslog/pkg/log"
)

var payload = []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)
var err error

func BenchmarkParsing(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err = syslog.Parse(payload)
	}
}

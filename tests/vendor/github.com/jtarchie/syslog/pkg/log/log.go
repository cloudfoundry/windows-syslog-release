package syslog

// go:generate ragel -e -G2 -Z parse.rl
//go:generate ragel -Z parse.rl

import "time"

type Property struct {
	Key   []byte
	Value []byte
}

type structureData struct {
	id         []byte
	properties []Property
}

func (s *structureData) ID() []byte {
	return s.id
}

func (s *structureData) Properties() []Property {
	return s.properties
}

type Log struct {
	version   int
	priority  int
	timestamp time.Time
	hostname  []byte
	appname   []byte
	procID    []byte
	msgID     []byte
	data      *structureData
	message   []byte
}

func (m *Log) Version() int {
	return m.version
}

func (m *Log) Severity() int {
	return m.priority & 7
}

func (m *Log) Facility() int {
	return m.priority >> 3
}

func (m *Log) Priority() int {
	return m.priority
}

func (m *Log) Timestamp() time.Time {
	return m.timestamp
}

func (m *Log) Hostname() []byte {
	return m.hostname
}

func (m *Log) Appname() []byte {
	return m.appname
}

func (m *Log) ProcID() []byte {
	return m.procID
}

func (m *Log) MsgID() []byte {
	return m.msgID
}

func (m *Log) StructureData() *structureData {
	return m.data
}

func (m *Log) Message() []byte {
	return m.message
}

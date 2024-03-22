package knock

import (
	"fmt"
	"net"
	"time"
)

type Event struct {
	t   time.Time
	msg string
}

func (e Event) String() string {
	return fmt.Sprintf("%s %s", e.t.Format("2006/01/02 15:04:05"), e.msg)
}

// ConnectionInfo 详细连接信息
type ConnectionInfo struct {
	events     []Event
	remoteAddr *net.TCPAddr
	localAddr  []*net.TCPAddr
	inter      net.Interface
	err        error
}

func (ci *ConnectionInfo) log(stage, msg string) {
	ci.events = append(ci.events, Event{
		t:   time.Now(),
		msg: "<" + stage + ">" + msg,
	})
}

// Snapshot 快照
type Snapshot struct {
	req  Request
	ci   *ConnectionInfo
	resp *Response
}

func (s *Snapshot) LocalAddr() []*net.TCPAddr {
	return s.ci.localAddr
}

func (s *Snapshot) RemoteAddr() *net.TCPAddr {
	return s.ci.remoteAddr
}

func (s *Snapshot) NetInterface() net.Interface {
	return s.ci.inter
}

func (s *Snapshot) Events() []Event {
	return s.ci.events
}

func (s *Snapshot) Request() Request {
	return s.req
}

func (s *Snapshot) Response() (*Response, error) {
	return s.resp, s.ci.err
}

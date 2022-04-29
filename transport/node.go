package transport

import (
	"github.com/mmatczuk/scylla-go-driver/frame"
	"go.uber.org/atomic"
)

type nodeStatus = atomic.Bool

const (
	statusDown = false
	statusUP   = true
)

type Node struct {
	hostID     frame.UUID
	addr       string
	datacenter string
	rack       string
	pool       *ConnPool
	status     nodeStatus
}

func (n *Node) Status() bool {
	return n.status.Load()
}

func (n *Node) setStatus(v bool) {
	n.status.Store(v)
}

func (n *Node) LeastBusyConn() *Conn {
	return n.pool.LeastBusyConn()
}

func (n *Node) Conn(token Token) *Conn {
	return n.pool.Conn(token)
}

type RingEntry struct {
	node  *Node
	token Token
}

func (r RingEntry) Less(i RingEntry) bool {
	return r.token < i.token
}

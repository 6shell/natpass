package handler

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jkstack/natpass/code/network"
	"github.com/lwch/logging"
)

type clients struct {
	sync.RWMutex
	parent *Handler
	id     string
	data   map[uint32]*client // idx => client
	idx    uint32
}

func newClients(parent *Handler, id string) *clients {
	logging.Info("new clients: %s", id)
	return &clients{
		parent: parent,
		id:     id,
		data:   make(map[uint32]*client),
	}
}

func (cs *clients) new(idx uint32, conn *network.Conn) *client {
	logging.Info("new client: %s-%d", cs.id, idx)
	cli := &client{
		parent:  cs,
		idx:     idx,
		conn:    conn,
		updated: time.Now(),
		links:   make(map[string]struct{}),
	}
	cs.Lock()
	cs.data[idx] = cli
	cs.Unlock()
	return cli
}

func (cs *clients) next() *client {
	list := make([]*client, 0, len(cs.data))
	cs.RLock()
	for _, cli := range cs.data {
		list = append(list, cli)
	}
	cs.RUnlock()
	if len(list) > 0 {
		idx := atomic.AddUint32(&cs.idx, 1)
		cli := list[int(idx)%len(list)]
		return cli
	}
	return nil
}

func (cs *clients) close(idx uint32) {
	cs.Lock()
	delete(cs.data, idx)
	cs.Unlock()
	if len(cs.data) == 0 {
		cs.parent.removeClients(cs.id)
	}
}

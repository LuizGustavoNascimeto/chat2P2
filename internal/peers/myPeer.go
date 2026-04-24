package peers

import (
	"sync"
)

type MyPeer struct {
	mu          sync.RWMutex
	Name        string
	InDiscovery bool
}

func NewMyPeer() *MyPeer {
	return &MyPeer{}
}

func (p *MyPeer) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Name = name
}

func (p *MyPeer) GetName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Name
}

func (p *MyPeer) SetInDiscovery(inDiscovery bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.InDiscovery = inDiscovery
}

func (p *MyPeer) IsInDiscovery() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.InDiscovery
}

func (p *MyPeer) Snapshot() MyPeer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return MyPeer{
		Name:        p.Name,
		InDiscovery: p.InDiscovery,
	}
}

package registry

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Gossip struct {
	registry *Registry
	peers    []string
	port     string
}

func NewGossip(reg *Registry, peers []string, port string) *Gossip {
	return &Gossip{
		registry: reg,
		peers:    peers,
		port:     port,
	}
}

func (g *Gossip) Start() {
	go func() {
		for {
			g.sync()
			time.Sleep(30 * time.Second)
		}
	}()
}

func (g *Gossip) sync() {
	for _, peer := range g.peers {
		data, err := g.fetch(peer)
		if err != nil {
			fmt.Printf("gossip: failed to reach %s: %v\n", peer, err)
			continue
		}
		if err := g.registry.FromJSON(data); err != nil {
			fmt.Printf("gossip: bad data from %s: %v\n", peer, err)
			continue
		}
		fmt.Printf("gossip: synced with %s\n", peer)
	}
}

func (g *Gossip) fetch(peer string) ([]byte, error) {
	resp, err := http.Get("http://" + peer + "/registry")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (g *Gossip) AddPeer(peer string) {
	g.peers = append(g.peers, peer)
}

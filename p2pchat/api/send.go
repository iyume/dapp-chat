package api

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
)

// Backend defines events trigger
type Backend struct {
	key                     *ecdsa.PrivateKey
	server                  *p2p.Server
	peers                   []*Peer // peer/rw matrix
	running                 bool
	emitP2PMessageEvent     chan *P2PMessageEvent
	emitChannelMessageEvent chan *ChannelMessageEvent

	stop <-chan int // TODO: interrupt notify and disconnect all peers
}

type BackendConfig struct {
	Key      *ecdsa.PrivateKey
	MaxPeers int
	NAT      nat.Interface // p2p.nat.Parse
	Host     string
	Port     int
	Locally  bool // restrict local request

	BootstrapNodes []*enode.Node // enode.MustParse

	// for test
	NetRestrict *netutil.Netlist // p2p.netutil.ParseNetlist
}

var localCIDRs = func() *netutil.Netlist {
	netlist, err := netutil.ParseNetlist("127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16")
	if err != nil {
		log.Panicln(err)
	}
	return netlist
}()

// start p2p server and backend
func StartBackend(config BackendConfig, stop <-chan int) {
	var backend = new(Backend)
	if config.Locally {
		config.NetRestrict = localCIDRs
	}
	server := &p2p.Server{
		Config: p2p.Config{
			PrivateKey:     config.Key,
			MaxPeers:       config.MaxPeers,
			NAT:            config.NAT,
			Protocols:      MakeProtocols(backend),
			ListenAddr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			NetRestrict:    config.NetRestrict,
			BootstrapNodes: config.BootstrapNodes,
		},
	}
	*backend = Backend{
		key:                     config.Key,
		server:                  server,
		emitP2PMessageEvent:     make(chan *P2PMessageEvent),
		emitChannelMessageEvent: make(chan *ChannelMessageEvent),
		stop:                    stop,
	}
	if err := server.Start(); err != nil {
		log.Panicln(err)
	}
	// srv.LocalNode().Node() ensure localnode exists. srv.Self() will create it.
	log.Println("Started P2P networking at", server.LocalNode().Node().URLv4())
	log.Println("Node ID:", server.LocalNode().ID())
	go backend.Run()
}

func (backend *Backend) addPeer(p *Peer) {
	backend.peers = append(backend.peers, p)
}

func (backend *Backend) findPeer(nodeID [32]byte) *Peer {
	// we could check protocol version in further
	for _, p := range backend.peers {
		pID := p.p.ID()
		if bytes.Equal(pID[:], nodeID[:]) && !p.closed {
			return p
		}
	}
	return nil
}

func truncateNodeID(nodeID [32]byte) string {
	return string(nodeID[:8]) + "..."
}

func stringToIDV4(s string) ([32]byte, error) {
	nodeIDHex, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	nodeID := [32]byte{}
	copy(nodeID[:], nodeIDHex)
	return nodeID, nil
}

// main loop of Backend to be goroutine
func (b *Backend) Run() {
	b.running = true
	for {
		select {
		case <-b.stop:
			for _, p := range b.server.Peers() {
				p.Disconnect(p2p.DiscRequested)
			}
			return
		case p2pevent := <-b.emitP2PMessageEvent:
			nodeID, err := stringToIDV4(p2pevent.NodeID)
			if err != nil {
				log.Println("node ID is not valid hex string")
				continue
			}
			peer := b.findPeer(nodeID)
			if peer == nil {
				log.Printf(
					"p2p connection to %s is not established\n", truncateNodeID(nodeID))
				continue
			}
			// the peer connection is secure enough, but we could use ECIES/ECDH
			// for futher security
			p2p.Send(peer.rw, P2PMessageEventMsg, p2pevent)
		case <-b.emitChannelMessageEvent:
			log.Println("send channel message is not supported")
		}
	}
}

func (b *Backend) SendP2PMessage(nodeID string, message Message) error {
	if !b.running {
		return errors.New("p2p server not started")
	}
	if len(nodeID) != 64 {
		return errors.New("node ID must be string of length 64")
	}
	event := MakeP2PMessageEvent(time.Now(), message, nodeID)
	log.Println("sending p2p event:", event)
	b.emitP2PMessageEvent <- &event
	return nil
}

func (b *Backend) SendChannelMessage() {
	if !b.running {
		log.Panicln("p2p server not started")
	}
}

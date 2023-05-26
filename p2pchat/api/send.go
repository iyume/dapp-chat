package api

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/iyume/dapp-chat/p2pchat/utils"
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
	MaxPeers int           // default to 50
	NAT      nat.Interface // p2p.nat.Parse
	Address  string        // default to "127.0.0.1:0"
	Locally  bool          // restrict local request

	BootstrapNodes []*enode.Node // enode.MustParse

	// for test
	NetRestrict *netutil.Netlist // p2p.netutil.ParseNetlist
}

var DefaultBackendConfig = BackendConfig{
	MaxPeers: 50,
	NAT:      nat.Any(),
	Address:  "127.0.0.1:0",
}

var localCIDRs = func() *netutil.Netlist {
	netlist, err := netutil.ParseNetlist("127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16")
	if err != nil {
		log.Panicln(err)
	}
	return netlist
}()

func NewBackend(config BackendConfig, stop <-chan int) *Backend {
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
			ListenAddr:     config.Address,
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
	return backend
}

func (b Backend) NodeID() [32]byte {
	ln := b.server.LocalNode()
	if ln == nil {
		log.Panicln("backend is not started")
	}
	return ln.ID()
}

// Start p2p server and backend in goroutine
func (b *Backend) Start() {
	if err := b.server.Start(); err != nil {
		log.Panicln(err)
	}
	// srv.LocalNode().Node() ensure localnode exists. srv.Self() will create it.
	log.Println("Started P2P networking at", b.server.LocalNode().Node().URLv4())
	log.Println("Node ID:", b.server.LocalNode().ID().String())
	go b.run()
}

func (b *Backend) addPeer(p *Peer) {
	b.peers = append(b.peers, p)
}

func (b *Backend) findPeer(nodeID [32]byte) *Peer {
	// we could check protocol version in further
	for _, p := range b.peers {
		if bytes.Equal(p.p.ID().Bytes(), nodeID[:]) && !p.closed {
			return p
		}
	}
	return nil
}

func truncateBytes(nodeID [32]byte) string {
	return string(nodeID[:8]) + "..."
}

// main loop of Backend to be goroutine
func (b *Backend) run() {
	b.running = true
	for {
		select {
		case <-b.stop:
			for _, p := range b.server.Peers() {
				p.Disconnect(p2p.DiscRequested)
			}
			return
		case p2pevent := <-b.emitP2PMessageEvent:
			nodeID, err := utils.ParseHexNodeID(p2pevent.NodeID)
			if err != nil {
				log.Println("node ID is not valid hex string")
				continue
			}
			peer := b.findPeer(nodeID)
			if peer == nil {
				log.Printf(
					"p2p connection to %s is not established\n", truncateBytes(nodeID))
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

// Backend APIs

type peersInfo struct {
	NodeID  [32]byte `json:"node_id"`
	Active  bool     `json:"active"`
	Version uint     `json:"version"`
}

func (b *Backend) PeersInfo() *[]peersInfo {
	infos := []peersInfo{}
	for _, p := range b.peers {
		infos = append(infos, peersInfo{
			NodeID:  p.p.ID(),
			Active:  !p.closed,
			Version: p.version,
		})
	}
	return &infos
}

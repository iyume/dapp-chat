package api

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/types"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

// Backend defines events trigger
type Backend struct {
	key    *ecdsa.PrivateKey
	server *p2p.Server
	peers  []*Peer // peer/rw matrix
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
	Locally:  true,
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

func NewBackend(config BackendConfig) *Backend {
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
		key:    config.Key,
		server: server,
	}
	return backend
}

func (b *Backend) NodeID() enode.ID {
	ln := b.server.LocalNode()
	if ln == nil {
		log.Panicln("backend is not started")
	}
	return ln.ID()
}

// SessionID returns p2p session ID for database, empty bytes if nodeID equals to self ID
func (b *Backend) SessionID(nodeID [32]byte) [32]byte {
	return utils.GetSessionID(b.NodeID(), nodeID)
}

func (b *Backend) Stop() {
	for _, p := range b.server.Peers() {
		p.Disconnect(p2p.DiscQuitting)
	}
}

// Start p2p server and backend in goroutine
func (b *Backend) Start() {
	if err := b.server.Start(); err != nil {
		log.Panicln(err)
	}
	// srv.LocalNode().Node() ensure localnode exists. srv.Self() will create it.
	log.Println("Started P2P networking at", b.server.LocalNode().Node().URLv4())
	log.Println("Node ID:", b.server.LocalNode().ID().String())
}

func (b *Backend) addPeer(p *Peer) {
	b.peers = append(b.peers, p)
}

func (b *Backend) findPeer(nodeID [32]byte) *Peer {
	// we could check protocol version in further (p.RunningCap)
	for _, p := range b.peers {
		if bytes.Equal(p.p.ID().Bytes(), nodeID[:]) && !p.closed {
			return p
		}
	}
	return nil
}

func truncateBytes(nodeID [32]byte) string {
	return hex.EncodeToString(nodeID[:4]) + "..."
}

func (b *Backend) SendP2PMessage(nodeID [32]byte, message types.Message) error {
	if message.Empty() {
		return errors.New("message is empty")
	}
	if bytes.Equal(b.NodeID().Bytes(), nodeID[:]) {
		return errors.New("cannot send message to self")
	}
	p := b.findPeer(nodeID)
	if p == nil {
		msg := fmt.Sprintf("connection to %s is not established", truncateBytes(nodeID))
		log.Println(msg)
		return errors.New(msg)
	}
	log.Printf("sent p2p message '%s' to node %s\n", message.ExtractPlaintext(), truncateBytes(nodeID))
	event := types.MakeP2PMessageEvent(message)
	if err := p.SendMessage(event); err != nil {
		return err
	}
	// Message is sent properly, but we don't know if message is properly received
	// TODO: fix it
	event.UserID = b.NodeID().String()
	db.AddP2PMessageEvent(utils.GetSessionID(b.NodeID(), nodeID), event)
	return nil
}

// Backend APIs

type peersInfo struct {
	NodeID  string `json:"node_id"`
	Active  bool   `json:"active"`
	Version uint   `json:"version"`
}

func (b *Backend) PeersInfo() *[]peersInfo {
	infos := []peersInfo{}
	for _, p := range b.peers {
		infos = append(infos, peersInfo{
			NodeID:  p.p.ID().String(),
			Active:  !p.closed,
			Version: p.version,
		})
	}
	return &infos
}

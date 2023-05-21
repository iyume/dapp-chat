package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/iyume/dapp-chat/p2pchat/api"
)

const maxBodySize = 1 * 1024 * 1024

// TODO: refactor config

var key = getGenerateKey()

var backend = api.NewBackend(key)

// p2pserver must be started before http server
var p2pserver = &p2p.Server{
	Config: p2p.Config{
		PrivateKey: key,
		MaxPeers:   50,
		NAT:        nil, // equals to 'none'
		Protocols:  api.MakeProtocols(backend),
		ListenAddr: ":0",
	},
}

var mux = http.NewServeMux()

var server = &http.Server{
	Handler:      mux,
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func flatQueryParams(r *http.Request) map[string]string {
	params := r.URL.Query()
	res := map[string]string{}
	for key, value := range params {
		res[key] = value[0]
	}
	return res
}

func readBody(r *http.Request) ([]byte, error) {
	bytes, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		return nil, err
	}
	if len(bytes) >= maxBodySize {
		return nil, errors.New("exceeded max body size")
	}
	r.Body.Close()
	return bytes, err
}

func init() {
	initServer()
	initP2PConfig()
}

func initP2PConfig() {
	// for test
	netrestrict, err := netutil.ParseNetlist("127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16")
	if err != nil {
		log.Fatalln(err)
	}
	p2pserver.Config.NetRestrict = netrestrict
	if bootnode != "" {
		p2pserver.Config.BootstrapNodes = []*enode.Node{enode.MustParse(bootnode)}
	}
}

func initServer() {
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})
	mux.HandleFunc("/send_p2p_message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		params := flatQueryParams(r)
		data := QuerySendP2PMessage{
			Nickname: params["nickname"],
			NodeID:   params["node_id"],
			Message:  params["message"],
		}
		if r.Method == "POST" {
			bytes, err := readBody(r)
			if err != nil {
				log.Println(err)
				return
			}
			err = json.Unmarshal(bytes, &data)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				log.Println(err)
				return
			}
		}
		if data.HasZeroField() {
			http.Error(w, "parameter not enough", http.StatusBadRequest)
			return
		}
		// convert to Message if string
		var message api.Message
		switch messageVal := data.Message.(type) {
		case string:
			message = api.Message{
				api.Segment{Type: "text", Data: api.TextSegment{Text: messageVal}},
			}
		case api.Message:
			message = messageVal
		}
		if data.NodeID == "" {
			// we couldn't use Nickname without database
			http.Error(w, "node ID must be provided", http.StatusBadRequest)
			return
		}
		backend.SendP2PMessage(data.NodeID, message)
		w.Write([]byte(`{"message_id":0}`))
	})
}
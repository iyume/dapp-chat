package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/types"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

const (
	defaultIPFSDataDir = "/.p2pchat"

	// /.p2pchat/self_id
	ipfsSelfIDFile = "self_id"

	// /.p2pchat/p2p_sessions
	ipfsP2PSessionsDir = "p2p_sessions"

	// /.p2pchat/p2p_sessions/xxx/index
	ipfsP2PSessionIndex = "index"
)

type Caller interface {
	// Call returns error should be handled by server
	Call(string, Getter) map[string]any
}

// may implemented by codegen

var actions = map[string]func(b *api.Backend, p Getter) map[string]any{
	"get_self_id": func(b *api.Backend, p Getter) map[string]any {
		return OK(b.SelfID().String())
	},
	"get_peers_info": func(b *api.Backend, p Getter) map[string]any {
		return OK(b.PeersInfo())
	},
	"get_friend_list": func(b *api.Backend, p Getter) map[string]any {
		resp := []map[string]any{}
		for nodeID, info := range db.GetFriends() {
			resp = append(resp, map[string]any{
				"node_id": hex.EncodeToString(nodeID[:]),
				"remark":  info.Remark,
			})
		}
		return OK(resp)
	},
	"add_friend": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id", "remark"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		remark := p.GetString("remark")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		db.AddFriend(nodeID, remark)
		return OK(nil)
	},
	"delete_friend": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		db.DeleteFriend(nodeID)
		return OK(nil)
	},
	"get_p2p_session": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		session := db.GetP2PSession(b.SessionID(nodeID))
		if session == nil {
			return Failed("internal error")
		} else {
			return OK(session)
		}
	},
	"send_p2p_message": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id", "message"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		messageIntf := p.Get("message")
		var message types.Message
		switch val := messageIntf.(type) {
		case string:
			message = types.PlaintextToMessage(val)
		case types.Message:
			message = val
		default:
			return Failed("parameter invalid")
		}
		if message.Empty() {
			return Failed("message is empty")
		}
		if err := b.SendP2PMessage(nodeID, message); err != nil {
			log.Println("error sending p2p message:", err)
			return Failed("internal error")
		}
		return OK(nil)
	},
	"upload_ipfs": func(b *api.Backend, p Getter) map[string]any {
		// upload messages that sent by me to IPFS
		if miss := p.Require("ipfs_addr"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		ipfsAddr := p.GetString("ipfs_addr")
		mfsDataDir := p.GetString("mfs_data_dir")
		if mfsDataDir == "" {
			mfsDataDir = defaultIPFSDataDir
		} else if !strings.HasPrefix(mfsDataDir, "/") {
			return Failed("parameter 'mfs_data_dir' must be startswith '/'")
		}
		key := p.GetString("key")
		if key == "" {
			key = "self"
		}

		sessions := db.GetP2PSessions()
		// message hash to encrypted message object
		// NOTE: json cannot marshal map[[32]byte]xxx (needs comparable), so use hex string
		outBoundMessages := map[string][]byte{}
		// target ID to list of message hash
		outBoundSessions := map[string][]string{}

		for _, session := range sessions {
			if session.Pubkey == nil {
				log.Println("session is not annotated with pubkey and nodeID")
				continue
			}
			pubkey, err := crypto.UnmarshalPubkey(session.Pubkey)
			if err != nil {
				log.Println("error unmarshal pubkey for", session.UserID)
				return Failed("internal error")
			}
			for _, event := range session.Events {
				if event.UserID == b.SelfID().String() {
					targetID := session.UserID.String()
					// NOTE: json marshal respects the order of struct and sort map key
					data, err := json.Marshal(event)
					if err != nil {
						log.Println("cannot marshal event", event)
						return Failed("internal error")
					}
					hashBytes := sha256.Sum256(data)
					hash := hex.EncodeToString(hashBytes[:])
					ciphertext, err := b.Encrypt(pubkey, data)
					if err != nil {
						log.Println("error encrypt data")
						return Failed("internal error")
					}
					outBoundMessages[hash] = ciphertext
					outBoundSessions[targetID] = append(outBoundSessions[targetID], hash)
				}
			}
		}

		log.Println("dialing with", ipfsAddr)
		sh := shell.NewShell(ipfsAddr)
		log.Println("cleaning IPFS data directory", mfsDataDir)
		sh.FilesRm(context.Background(), mfsDataDir, true) // the official recursive option?
		// Use String for readable
		cid, err := sh.Add(strings.NewReader(b.SelfID().String()))
		if err != nil {
			return Failed(err.Error())
		}
		log.Printf("uploaded self_id at /ipfs/%s\n", cid)
		if err := sh.FilesCp(context.Background(), "/ipfs/"+cid, path.Join(mfsDataDir, ipfsSelfIDFile),
			shell.FilesCp.Parents(true),
		); err != nil {
			return Failed(err.Error())
		}
		sessionsDir := path.Join(mfsDataDir, ipfsP2PSessionsDir)
		if err := sh.FilesMkdir(context.Background(), sessionsDir); err != nil {
			return Failed(err.Error())
		}
		for targetID, msgHashes := range outBoundSessions {
			sessionDir := path.Join(sessionsDir, targetID)
			log.Println("doing mkdir", sessionDir)
			if err := sh.FilesMkdir(context.Background(), sessionDir); err != nil {
				return Failed(err.Error())
			}
			// Since IPFS not yet supports list directory in gateway. Here indexes all hashes
			// See: https://github.com/ipfs/kubo/issues/7552
			indexBuf := &bytes.Buffer{}
			for _, msgHash := range msgHashes {
				cid, err := sh.Add(bytes.NewReader(outBoundMessages[msgHash]))
				if err != nil {
					return Failed(err.Error())
				}
				log.Printf("uploaded message at /ipfs/%s\n", cid)
				msgFile := path.Join(sessionDir, msgHash)
				if err := sh.FilesCp(context.Background(), "/ipfs/"+cid, msgFile); err != nil {
					return Failed(err.Error())
				}
				log.Printf("mapped /ipfs/%s to MFS %s\n", cid, msgFile)
				indexBuf.WriteString(msgHash)
				indexBuf.WriteByte('\n')
			}
			cid, err := sh.Add(indexBuf)
			if err != nil {
				return Failed(err.Error())
			}
			if err := sh.FilesCp(context.Background(),
				"/ipfs/"+cid, path.Join(sessionDir, ipfsP2PSessionIndex),
			); err != nil {
				return Failed(err.Error())
			}
		}

		// publish to IPNS
		log.Println("publishing with key", key)
		stat, err := sh.FilesStat(context.Background(), mfsDataDir)
		if err != nil {
			return Failed(err.Error())
		}
		publishResp, err := sh.PublishWithDetails(stat.Hash, key, 0, 0, true)
		if err != nil {
			return Failed(err.Error())
		}
		log.Println("published", *publishResp)

		return OK(publishResp)
	},
	"verify_ipfs": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("ipfs_addr"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		return OK(nil)
	},
}

type caller struct {
	backend *api.Backend
}

func (c caller) Call(action string, p Getter) map[string]any {
	handler, ok := actions[action]
	if !ok {
		return nil
	}
	return handler(c.backend, p)
}

func NewCaller(b *api.Backend) Caller {
	return caller{backend: b}
}

func OK(data any) map[string]any {
	return map[string]any{"retcode": 0, "data": data}
}

func Failed(reason string) map[string]any {
	return map[string]any{"retcode": 1, "reason": reason}
}

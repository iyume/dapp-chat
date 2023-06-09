package db

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/iyume/dapp-chat/p2pchat/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Abstract database operations for multi database support
// type Database interface {
// }

type friendInfo struct {
	Remark string `json:"remark"`
}

type p2pSession struct {
	UserID enode.ID                `json:"_user_id"`
	Pubkey []byte                  `json:"_pubkey"`
	Events []types.P2PMessageEvent `json:"events"`
}

// initial with allocated content
func newP2PSession() *p2pSession {
	return &p2pSession{Events: make([]types.P2PMessageEvent, 0)}
}

// leveldb

type Database struct {
	*leveldb.DB

	inited   bool
	initOnce sync.Once
}

var _ldb = new(Database)
var ErrDBNotInit = errors.New("database not initialize")

func getDatabase() *Database {
	if !_ldb.inited {
		panic(ErrDBNotInit)
	}
	return _ldb
}

func newPersistentDB(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return db
}

func doInit(path string, localNodeID [32]byte) {
	ldb := newPersistentDB(path)
	exists, err := ldb.Has(localKey, nil)
	if err != nil {
		panic(err)
	}
	if exists {
		data, err := ldb.Get(localKey, nil)
		if err != nil {
			panic(err)
		}
		if !(bytes.Equal(localNodeID[:], data)) {
			panic(fmt.Sprintf(
				"database has been already initialized with Node ID %x", data,
			))
		}
	} else {
		if err := ldb.Put(localKey, localNodeID[:], nil); err != nil {
			panic(err)
		}
	}
	_ldb.DB = ldb
	_ldb.inited = true
	log.Printf("leveldb initialized at %s with Node ID %x\n", path, localNodeID)
}

// Register leveldb database
func Init(path string, localNodeID [32]byte) {
	_ldb.initOnce.Do(func() { doInit(path, localNodeID) })
}

// This returns non-nil value.
func GetFriendIDs() [][32]byte {
	db := getDatabase()
	fIDs := [][32]byte{}
	iter := db.NewIterator(util.BytesPrefix(friendPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		fIDs = append(fIDs, [32]byte(iter.Key()[len(friendPrefix):]))
	}
	if err := iter.Error(); err != nil {
		log.Println(err)
		return [][32]byte{}
	}
	return fIDs
}

func HasFriend(nodeID [32]byte) bool {
	db := getDatabase()
	// ids := GetFriendIDs()
	// return slices.Contains(*ids, nodeID)
	exist, err := db.Has(friendKey(nodeID), nil)
	if err != nil {
		log.Println(err)
		return false
	}
	return exist
}

// Get all friends info. Returns nil if any error occurs
func GetFriends() map[[32]byte]friendInfo {
	db := getDatabase()
	friends := map[[32]byte]friendInfo{}
	iter := db.NewIterator(util.BytesPrefix(friendPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		info := friendInfo{}
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			log.Println(errors.Join(err,
				fmt.Errorf("cannot unmarshal value on key '%s'", iter.Key())))
			return nil
		}
		friends[[32]byte(iter.Key()[len(friendPrefix):])] = info
	}
	if err := iter.Error(); err != nil {
		log.Println(err)
		return nil
	}
	return friends
}

// Add new friend with associated info
func AddFriend(nodeID [32]byte, remark string) {
	db := getDatabase()
	data, err := json.Marshal(friendInfo{Remark: remark})
	if err != nil {
		log.Println(err)
		return
	}
	// shall we check friend existence?
	if err := db.Put(friendKey(nodeID), data, nil); err != nil {
		log.Println(err)
		return
	}
}

func DeleteFriend(nodeID [32]byte) {
	db := getDatabase()
	if err := db.Delete(friendKey(nodeID), nil); err != nil {
		log.Println(err)
		return
	}
}

// Get specific p2p session. Returns nil if any error occurs
func GetP2PSession(sessionID [32]byte) *p2pSession {
	db := getDatabase()
	data, err := db.Get(p2pSessionKey(sessionID), nil)
	if err != nil {
		log.Println(err)
		return newP2PSession()
	}
	session := new(p2pSession)
	if err := json.Unmarshal(data, session); err != nil {
		log.Println(err)
		return newP2PSession()
	}
	return session
}

// Get all p2p sessions. Returns nil if any error occurs
func GetP2PSessions() []*p2pSession {
	db := getDatabase()
	sessions := []*p2pSession{}
	iter := db.NewIterator(util.BytesPrefix(p2pSessionPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		session := new(p2pSession)
		if err := json.Unmarshal(iter.Value(), session); err != nil {
			log.Println(errors.Join(err, fmt.Errorf("cannot unmarshal value on key '%s'", iter.Key())))
			return nil
		}
		sessions = append(sessions, session)
	}
	if err := iter.Error(); err != nil {
		log.Println(err)
		return nil
	}
	return sessions
}

func PutP2PSession(sessionID [32]byte, session *p2pSession) {
	db := getDatabase()
	data, err := json.Marshal(session)
	if err != nil {
		log.Println(err)
		return
	}
	if err := db.Put(p2pSessionKey(sessionID), data, nil); err != nil {
		log.Println(err)
		return
	}
}

func NotateP2PSession(sessionID [32]byte, userID enode.ID, pubkey *ecdsa.PublicKey) {
	session := GetP2PSession(sessionID)
	if session == nil {
		log.Println("cannot get p2p session")
		return
	}
	if session.Pubkey == nil {
		session.Pubkey = crypto.FromECDSAPub(pubkey)
		session.UserID = userID
		PutP2PSession(sessionID, session)
	}
}

// Add appends p2p message event into session, creates session if not exists.
func AddP2PMessageEvent(sessionID [32]byte, event types.P2PMessageEvent) {
	if event.UserID == "" {
		log.Println("missing UserID field")
		return
	}
	session := GetP2PSession(sessionID)
	if session == nil {
		log.Println("cannot get p2p session")
		return
	}
	session.Events = append(session.Events, event)
	PutP2PSession(sessionID, session)
}

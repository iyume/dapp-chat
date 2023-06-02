package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
	Events []types.P2PMessageEvent `json:"events"`
}

// initial with allocated content
func newP2PSession() *p2pSession {
	return &p2pSession{Events: make([]types.P2PMessageEvent, 0)}
}

// leveldb

type Database struct {
	*leveldb.DB

	inited bool
}

var db = new(Database)
var ErrDBNotInit = errors.New("database not initialized")

func newPersistentDB(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Panicln(err)
	}
	return db
}

// Register database
func Init(path string, localNodeID [32]byte) error {
	*db = Database{DB: newPersistentDB(path)}
	exists, err := db.Has(localKey, nil)
	if err != nil {
		return err
	}
	if exists {
		data, err := db.Get(localKey, nil)
		if err != nil {
			return err
		}
		if !(bytes.Equal(localNodeID[:], data)) {
			return fmt.Errorf("database has been already initialized with Node ID %x", data)
		}
	} else {
		if err := db.Put(localKey, localNodeID[:], nil); err != nil {
			return err
		}
	}
	db.inited = true
	return nil
}

func GetFriendIDs() [][32]byte {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	fIDs := [][32]byte{}
	iter := db.NewIterator(util.BytesPrefix(friendPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		fIDs = append(fIDs, [32]byte(iter.Key()[len(friendPrefix):]))
	}
	if err := iter.Error(); err != nil {
		log.Println(err)
		return nil
	}
	return fIDs
}

func HasFriend(nodeID [32]byte) bool {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	// ids := GetFriendIDs()
	// return slices.Contains(*ids, nodeID)
	exist, err := db.Has(friendKey(nodeID), nil)
	if err != nil {
		log.Println(err)
		return false
	}
	return exist
}

func GetFriends() map[[32]byte]friendInfo {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
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
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
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
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	if err := db.Delete(friendKey(nodeID), nil); err != nil {
		log.Println(err)
		return
	}
}

// This returns non-nil value
func GetP2PSession(sessionID [32]byte) *p2pSession {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
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

func PutP2PSession(sessionID [32]byte, session *p2pSession) {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
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

// Add appends p2p message event into session, creates session if not exists.
func AddP2PMessageEvent(sessionID [32]byte, event types.P2PMessageEvent) {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	if event.UserID == "" {
		log.Println("missing UserID field")
		return
	}
	session := GetP2PSession(sessionID)
	session.Events = append(session.Events, event)
	PutP2PSession(sessionID, session)
}

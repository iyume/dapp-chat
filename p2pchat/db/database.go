package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Abstract database operations for multi database support
// type Database interface {
// }

type friendInfo struct {
	Remark string
}

type p2pSession struct {
	Messages []api.P2PMessageEvent
}

// leveldb

type Database struct {
	*leveldb.DB

	inited bool
}

var db *Database
var ErrDBNotInit = errors.New("database not initialized")

func newPersistentDB(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

// Register database
func Init(path string, localNodeID [32]byte) error {
	db = &Database{DB: newPersistentDB(path)}
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

func GetFriendIDs() *[][32]byte {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	friends := [][32]byte{}
	iter := db.NewIterator(util.BytesPrefix(friendPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		friends = append(friends, [32]byte(iter.Key()[len(friendPrefix):]))
	}
	if err := iter.Error(); err != nil {
		log.Panicln(err)
	}
	return &friends
}

func HasFriend(nodeID [32]byte) bool {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	// ids := GetFriendIDs()
	// return slices.Contains(*ids, nodeID)
	exist, err := db.Has(friendKey(nodeID), nil)
	if err != nil {
		log.Panicln(err)
	}
	return exist
}

func GetFriends() *map[[32]byte]friendInfo {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	friends := map[[32]byte]friendInfo{}
	iter := db.NewIterator(util.BytesPrefix(friendPrefix), nil)
	defer iter.Release()
	for iter.Next() {
		info := friendInfo{}
		if err := json.Unmarshal(iter.Value(), &info); err != nil {
			log.Panicln(errors.Join(err,
				fmt.Errorf("cannot unmarshal value on key '%s'", iter.Key())))
		}
		friends[[32]byte(iter.Key()[len(friendPrefix):])] = info
	}
	if err := iter.Error(); err != nil {
		log.Panicln(err)
	}
	return &friends
}

// Add new friend with associated info
func AddFriend(nodeID [32]byte, remark string) {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	data, err := json.Marshal(friendInfo{Remark: remark})
	if err != nil {
		log.Panicln(err)
	}
	// shall we check friend existence?
	if err := db.Put(friendKey(nodeID), data, nil); err != nil {
		log.Panicln(err)
	}
}

func DeleteFriend(nodeID [32]byte) {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	if err := db.Delete(friendKey(nodeID), nil); err != nil {
		log.Panicln(err)
	}
}

func GetP2PSession() (*p2pSession, error) {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	return nil, nil
}

func AddP2PMessage(sessionID [32]byte, message api.P2PMessageEvent) error {
	if !db.inited {
		log.Panicln(ErrDBNotInit)
	}
	return nil
}

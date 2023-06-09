package db

import (
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

const testDBPrefix = "__test_d180b4f3ac4509c5_"

var dbIDCounter uint32

func getDBPath() string {
	id := int(atomic.AddUint32(&dbIDCounter, 1))
	return testDBPrefix + strconv.Itoa(id)
}

func closeCleanly(db *leveldb.DB, dbPath string) {
	if err := db.Close(); err != nil {
		log.Panicln(err)
	}
	if err := os.RemoveAll(dbPath); err != nil {
		log.Panicln(err)
	}
}

type dbdata struct {
	key   []byte
	value []byte
}

func dumpDB(db *leveldb.DB) []dbdata {
	res := []dbdata{}
	for iter := db.NewIterator(nil, nil); iter.Next(); {
		res = append(res, dbdata{key: iter.Key(), value: iter.Value()})
	}
	return res
}

func TestInit(t *testing.T) {
	assert := assert.New(t)

	// panics not inited
	assert.Panics(func() { GetFriendIDs() })
	nodeID := [32]byte([]byte(
		"0076db4fee435414c8897271d126f0b356a5087e43e3cb5df12df73c482a6a2a"))
	dbPath := getDBPath()
	doInit(dbPath, nodeID)
	defer closeCleanly(getDatabase().DB, dbPath)
	assert.ElementsMatch([]dbdata{{key: localKey, value: nodeID[:]}}, dumpDB(getDatabase().DB))
}

func TestFriendCRUD(t *testing.T) {
	assert := assert.New(t)

	f1ID, f1Remark := [32]byte([]byte("ca9897c18db6a38d7a417c42380837e9426ff3171664a612e35c7ea15b70fb9f")), "friend1"
	f2ID, f2Remark := [32]byte([]byte("309844745a5d419c24d7ebd775bc5bc6b7791eaf45a393d86cacca5d489e22e4")), "friend2"
	dbPath := getDBPath()
	doInit(dbPath, [32]byte{})
	defer closeCleanly(getDatabase().DB, dbPath)
	assert.Equal(map[[32]byte]friendInfo{}, GetFriends())
	AddFriend(f1ID, f1Remark)
	assert.Equal(map[[32]byte]friendInfo{f1ID: {Remark: f1Remark}}, GetFriends())
	AddFriend(f2ID, f2Remark)
	assert.Equal(map[[32]byte]friendInfo{f1ID: {Remark: f1Remark}, f2ID: {Remark: f2Remark}}, GetFriends())
	AddFriend(f1ID, "friend1-renamed") // this is not proper behaviour
	assert.Equal(map[[32]byte]friendInfo{f1ID: {Remark: "friend1-renamed"}, f2ID: {Remark: f2Remark}}, GetFriends())
	assert.ElementsMatch([][32]byte{f1ID, f2ID}, GetFriendIDs())
	assert.False(HasFriend([32]byte{}))
	assert.True(HasFriend(f1ID))
	DeleteFriend(f1ID)
	assert.Equal(map[[32]byte]friendInfo{f2ID: {Remark: f2Remark}}, GetFriends())
}

package merkle

import (
	"bytes"
	"encoding/binary"
	"io"

	cmn "github.com/kardiachain/go-kardiamain/lib/common"
)

// Merkle tree from a map.
// Leaves are `hash(key) | hash(value)`.
// Leaves are sorted before Merkle hashing.
type simpleMap struct {
	kvs    cmn.KVPairs
	sorted bool
}

func newSimpleMap() *simpleMap {
	return &simpleMap{
		kvs:    nil,
		sorted: false,
	}
}

// Set creates a kv pair of the key and the hash of the value,
// and then appends it to simpleMap's kv pairs.
func (sm *simpleMap) Set(key string, value []byte) {
	sm.sorted = false

	// The value is hashed, so you can
	// check for equality with a cached value (say)
	// and make a determination to fetch or not.
	vhash := Sum(value)

	sm.kvs = append(sm.kvs, cmn.KVPair{
		Key:   []byte(key),
		Value: vhash,
	})
}

// Hash Merkle root hash of items sorted by key
// (UNSTABLE: and by value too if duplicate key).
func (sm *simpleMap) Hash() []byte {
	sm.Sort()
	return hashKVPairs(sm.kvs)
}

func (sm *simpleMap) Sort() {
	if sm.sorted {
		return
	}
	sm.kvs.Sort()
	sm.sorted = true
}

// Returns a copy of sorted KVPairs.
// NOTE these contain the hashed key and value.
func (sm *simpleMap) KVPairs() cmn.KVPairs {
	sm.Sort()
	kvs := make(cmn.KVPairs, len(sm.kvs))
	copy(kvs, sm.kvs)
	return kvs
}

//----------------------------------------

// A local extension to KVPair that can be hashed.
// Key and value are length prefixed and concatenated,
// then hashed.
type KVPair cmn.KVPair

// Bytes returns key || value, with both the
// key and value length prefixed.
func (kv KVPair) Bytes() []byte {
	var b bytes.Buffer
	err := EncodeByteSlice(&b, kv.Key)
	if err != nil {
		panic(err)
	}
	err = EncodeByteSlice(&b, kv.Value)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func hashKVPairs(kvs cmn.KVPairs) []byte {
	kvsH := make([][]byte, len(kvs))
	for i, kvp := range kvs {
		kvsH[i] = KVPair(kvp).Bytes()
	}
	return SimpleHashFromByteSlices(kvsH)
}

func EncodeByteSlice(w io.Writer, bz []byte) (err error) {
	err = EncodeUvarint(w, uint64(len(bz)))
	if err != nil {
		return
	}
	_, err = w.Write(bz)
	return
}

// EncodeUvarint is used to encode golang's int, int32, int64 by default. unless specified differently by the
// `binary:"fixed32"`, `binary:"fixed64"`, or `binary:"zigzag32"` `binary:"zigzag64"` tags.
// It matches protobufs varint encoding.
func EncodeUvarint(w io.Writer, u uint64) (err error) {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], u)
	_, err = w.Write(buf[0:n])
	return
}

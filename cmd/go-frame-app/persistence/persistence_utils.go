package persistence

import (
	"encoding/binary"
	bolt "go.etcd.io/bbolt"
)

func isBucketEmpty(bucket *bolt.Bucket) bool {
	first, value := bucket.Cursor().First()
	return first == nil && value == nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func i64tob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

package farmhash

import (
	"reflect"
	"unsafe"

	farm "github.com/dgryski/go-farm"
)

// Hash32 returns 32 bits hash of a string.
func Hash32(s string) uint32 {
	// Using unsafe here to avoid escaping to heap / copying bytes on heap.
	// similar to reinterpret_cast in c++
	var b []byte
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
	return farm.Fingerprint32(b)
}

// Hash64 returns 64 bits hash of a string.
func Hash64(s string) uint64 {
	// Using unsafe here to avoid escaping to heap / copying bytes on heap.
	// similar to reinterpret_cast in c++
	var b []byte
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
	return farm.Fingerprint64(b)
}

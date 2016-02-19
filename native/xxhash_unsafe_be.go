// +build !safe
// +build !appengine
// +build be

package xxhash

import (
	"reflect"
	"unsafe"
)

var (
	dummy = [2]byte{1, 0}
	isBig = *(*int16)(unsafe.Pointer(&dummy[0])) != 1
)

// Backend returns the current version of xxhash being used.
const Backend = "GoUnsafeBigEndian"

type byteReader struct {
	p unsafe.Pointer
}

func newbyteReader(b []byte) byteReader {
	return byteReader{unsafe.Pointer(&b[0])}
}

func (br byteReader) Uint32(i int) (u uint32) {
	u = *(*uint32)(unsafe.Pointer(uintptr(br.p) + uintptr(i)))
	if isBig {
		u = swap32be(u)
	}
	return
}

func (br byteReader) Uint64(i int) (u uint64) {
	u = *(*uint64)(unsafe.Pointer(uintptr(br.p) + uintptr(i)))
	if isBig {
		u = swap64be(u)
	}
	return
}

func (br byteReader) Byte(i int) byte {
	return *(*byte)(unsafe.Pointer(uintptr(br.p) + uintptr(i)))
}

func swap32be(x uint32) uint32 {
	return ((x << 24) & 0xff000000) |
		((x << 8) & 0x00ff0000) |
		((x >> 8) & 0x0000ff00) |
		((x >> 24) & 0x000000ff)
}

func swap64be(x uint64) uint64 {
	x = (0xff00ff00ff00ff & (x >> 8)) | ((0xff00ff00ff00ff & x) << 8)
	x = (0xffff0000ffff & (x >> 16)) | ((0xffff0000ffff & x) << 16)
	return (0xffffffff & (x >> 32)) | ((0xffffffff & x) << 32)
}

func ChecksumString32(s string) uint32 {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	ss.Cap = ss.Len
	return Checksum32(*(*[]byte)(unsafe.Pointer(ss)))
}

func ChecksumString32S(s string, seed uint32) uint32 {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	ss.Cap = ss.Len
	return Checksum32S(*(*[]byte)(unsafe.Pointer(ss)), seed)
}

func ChecksumString64(s string) uint64 {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	ss.Cap = ss.Len
	return Checksum64(*(*[]byte)(unsafe.Pointer(ss)))
}

func ChecksumString64S(s string, seed uint64) uint64 {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	ss.Cap = ss.Len
	return Checksum64S(*(*[]byte)(unsafe.Pointer(ss)), seed)
}

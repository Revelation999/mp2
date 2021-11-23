package main

import (
	"bytes"
	"math"
	"unsafe"
)

func HeaderToByteSlice(header BlockHeader) []byte {
	var slice []byte
	slice = append(slice, IntToByteSlice(header.version)...)
	slice = append(slice, header.prevBlockHashPointer.hash[:]...)
	slice = append(slice, Int64ToByteSlice(int64(uintptr(unsafe.Pointer(header.prevBlockHashPointer.ptr))))...)
	slice = append(slice, header.merkleRootHashFiller...)
	slice = append(slice, IntToByteSlice(header.time)...)
	slice = append(slice, header.bits[:]...)
	slice = append(slice, IntToByteSlice(header.nonce)...)
	return slice
}

func IntToByteSlice(num int) []byte {
	var slice []byte
	if num == 0 {
		return append(slice, 0)
	}
	for true {
		if num > 0 {
			slice = append([]byte{byte(num % 256)}, slice...)
			num /= 256
		} else {
			break
		}
	}
	return slice
}

func Int64ToByteSlice(num int64) []byte {
	var slice []byte
	if num == 0 {
		return append(slice, 0)
	}
	for true {
		if num > 0 {
			slice = append([]byte{byte(num % 256)}, slice...)
			num /= 256
		} else {
			break
		}
	}
	return slice
}

func Compare(a, b []byte) int {
	for i := 0; i < int(math.Abs(float64(len(a)-len(b)))); i++ {
		if len(a) < len(b) {
			a = append([]byte{0}, a...)
		} else {
			b = append([]byte{0}, b...)
		}
	}
	return bytes.Compare(a, b)
}

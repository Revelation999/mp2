package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"time"
	"unsafe"
)

type BlockHeader struct {
	version int
	//hashPrevBlockHeader [32]byte
	prevBlockHashPointer HashPointer
	merkleRootHashFiller []byte
	time                 int
	bits                 [32]byte
	nonce                int
}

type HashPointer struct {
	hash [32]byte
	ptr  *Block
}

type Block struct {
	blockHeader BlockHeader
	transaction string
}

type Message struct {
	identity string
	nonce    int
	block    Block
}

type Logger struct {
	block         Block
	currBlockHash [32]byte
	mailbox       *chan Message
}

func newBlock(nonce int, provider string, prevBlock Block) Block {
	var prevBlockHashPointer HashPointer
	prevBlockHashPointer.hash = sha256.Sum256(HeaderToByteSlice(prevBlock.blockHeader)) // should we include the transaction?
	prevBlockHashPointer.ptr = &prevBlock
	var newHeader BlockHeader
	newHeader.version = prevBlock.blockHeader.version
	newHeader.prevBlockHashPointer = prevBlockHashPointer
	newHeader.merkleRootHashFiller = []byte{0}
	newHeader.time = int(time.Now().Unix())
	newHeader.bits = prevBlock.blockHeader.bits
	newHeader.nonce = nonce
	transaction := "Coin given to " + provider
	return Block{newHeader, transaction}
}

func (l *Logger) UpdateBlock(nonce int, provider string, miners []Miner) {
	l.block = newBlock(nonce, provider, l.block)
	fmt.Print("New block available: ")
	fmt.Println(l.block)
	l.currBlockHash = sha256.Sum256(HeaderToByteSlice(l.block.blockHeader))
	fmt.Print("Hash of current block: ")
	fmt.Println(l.currBlockHash)
	for i := 0; i < len(miners); i++ {
		*miners[i].mailbox <- l.block
	}
	for len(*l.mailbox) > 0 {
		<-*l.mailbox
	}
}

func (l Logger) CheckNonce(nonce int) bool {
	hashOutput := sha256.Sum256(append(l.currBlockHash[:], IntToByteSlice(nonce)...))
	return Compare(hashOutput[:], l.block.blockHeader.bits[:]) < 0
}

func (l Logger) ListenForUpdate(miners []Miner) {
	fmt.Println("Logger initiated.")
	for true {
		select {
		case msg := <-*l.mailbox:
			if l.CheckNonce(msg.nonce) {
				fmt.Print("Miner " + msg.identity + " has done it. It has found the value of ")
				fmt.Println(msg.nonce)
				fmt.Print("which generated a hash value of ")
				fmt.Println(sha256.Sum256(append(l.currBlockHash[:], IntToByteSlice(msg.nonce)...)))
				l.UpdateBlock(msg.nonce, msg.identity, miners)
			}
		default:
			continue
		}
		if time.Since(start).Nanoseconds() > TimeLimit {
			fmt.Println("Logger terminated at the 5 min mark.")
			break
		}
	}
}

func HeaderToByteSlice(header BlockHeader) []byte {
	var slice []byte
	slice = append(slice, IntToByteSlice(header.version)...)
	slice = append(slice, header.prevBlockHashPointer.hash[:]...)
	slice = append(slice, Int64ToByteSlice(int64(uintptr(unsafe.Pointer(header.prevBlockHashPointer.ptr))))...)
	//do we really need the pointer in the hash?
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
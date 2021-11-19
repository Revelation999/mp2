package main

import (
	"bytes"
	"crypto/sha256"
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

func (l Logger) UpdateBlock(nonce int, provider string, miners []Miner) {
	l.block = newBlock(nonce, provider, l.block)
	l.currBlockHash = sha256.Sum256(HeaderToByteSlice(l.block.blockHeader))
	for i := 0; i < len(miners); i++ {
		*miners[i].mailbox <- l.block
	}
}

func (l Logger) CheckNonce(nonce int) bool {
	hashOutput := sha256.Sum256(append(l.currBlockHash[:], IntToByteSlice(nonce)...))
	return bytes.Compare(hashOutput[:], l.block.blockHeader.bits[:]) < 0
}

func (l Logger) ListenForUpdate(miners []Miner) {
	for true {
		select {
		case msg := <-*l.mailbox:
			if l.CheckNonce(msg.nonce) {
				l.UpdateBlock(msg.nonce, msg.identity, miners)
			}
		default:
			continue
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

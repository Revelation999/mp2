package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Logger struct {
	block         Block
	currBlockHash [32]byte
	mailbox       *chan Message
}

func NewBlock(nonce int, provider string, prevBlock Block) Block {
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
	l.block = NewBlock(nonce, provider, l.block)
	puzzlesolvedin := time.Since(newpuzzle).Seconds()
	newpuzzle = time.Now()
	fmt.Println("It took ", puzzlesolvedin, " seconds to solve the puzzle.")
	fmt.Println("The following block has been added to the blockchain: ")
	PrettyPrintBlock(&l.block)
	l.currBlockHash = sha256.Sum256(HeaderToByteSlice(l.block.blockHeader))
	fmt.Println("\nMiners should now attempt to solve the puzzle given the following updated hash value: ")
	fmt.Println(l.currBlockHash, "\n")
	for i := 0; i < len(miners); i++ {
		*miners[i].mailbox <- l.block
	}
}

func (l Logger) CheckNonce(nonce int) bool {
	hashOutput := sha256.Sum256(append(l.currBlockHash[:], IntToByteSlice(nonce)...))
	return Compare(hashOutput[:], l.block.blockHeader.bits[:]) < 0
}

func (l Logger) ListenForUpdate(miners []Miner) {
	fmt.Println("Logger has begun listening for new updates.")
	for true {
		select {
		case msg := <-*l.mailbox:
			if l.CheckNonce(msg.nonce) {
				fmt.Println("\n\n! ! ! ! A NEW BLOCK HAS BEEN FOUND BY MINER ", msg.identity, " ! ! ! !")
				fmt.Println("Miner " + msg.identity + " solved the puzzle with a nonce value of", msg.nonce, ".")
				fmt.Println("This nonce generated a hash value of ", sha256.Sum256(append(l.currBlockHash[:], IntToByteSlice(msg.nonce)...)))
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

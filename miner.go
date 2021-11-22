package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Miner struct {
	identity      string
	block         Block
	mailbox       *chan Block
	currBlockHash [32]byte
}

func (m Miner) Mine(l Logger) {
	fmt.Println("Miner " + m.identity + " has begun mining.")
	i := 1
	for true {
		if m.HasUpdate() {
			fmt.Print("Miner " + m.identity + " has received a new puzzle.\n")
			i = 1
		}
		hashOutput := sha256.Sum256(append(m.currBlockHash[:], IntToByteSlice(i)...))
		//fmt.Println(hashOutput)
		//fmt.Println(m.block.blockHeader.bits)
		if Compare(hashOutput[:], m.block.blockHeader.bits[:]) < 0 {
			m.block = newBlock(i, m.identity, l.block)
			m.currBlockHash = sha256.Sum256(HeaderToByteSlice(l.block.blockHeader))
			*l.mailbox <- Message{m.identity, i, m.block}
		}
		i++
		if time.Since(start).Nanoseconds() > TimeLimit {
			fmt.Println("Miner " + m.identity + " terminated at the 5 min mark.")
			break
		}
	}
}

func (m *Miner) HasUpdate() bool {
	select {
	case block := <-*m.mailbox:
		m.block = block
		m.currBlockHash = sha256.Sum256(HeaderToByteSlice(m.block.blockHeader))
		return true
	default:
		return false
	}
}

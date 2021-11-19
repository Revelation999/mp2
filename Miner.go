package main

import (
	"bytes"
	"crypto/sha256"
)

type Miner struct {
	identity      string
	block         Block
	mailbox       *chan Block
	currBlockHash [32]byte
}

func (m Miner) Mine(l Logger) {
	i := 1
	for true {
		hashOutput := sha256.Sum256(append(m.currBlockHash[:], IntToByteSlice(i)...))
		if bytes.Compare(hashOutput[:], m.block.blockHeader.bits[:]) < 0 {
			m.block = newBlock(i, m.identity, l.block)
			m.currBlockHash = sha256.Sum256(HeaderToByteSlice(l.block.blockHeader))
			*l.mailbox <- Message{m.identity, i, m.block}
		}
		m.CheckUpdate()
	}
}

func (m Miner) CheckUpdate() {
	select {
	case block := <-*m.mailbox:
		m.block = block
	default:
		break
	}
}

package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

var logger Logger
var miners []Miner

const TimeLimit = 300_000_000_000

var start time.Time

func Genesis(emptyBlock *Block) Block {
	var hash [32]byte
	for i := 0; i < 32; i++ {
		hash[i] = byte(rand.Intn(256))
	}
	var difficulty [32]byte
	for i := 0; i < 32; i++ {
		if i != 2 {
			difficulty[i] = 0
		} else {
			difficulty[i] = 1
		}
	}
	hashPointer := HashPointer{hash, emptyBlock}
	newHeader := BlockHeader{1,
		hashPointer,
		[]byte{0},
		int(time.Now().Unix()),
		difficulty,
		0,
	}
	return Block{newHeader, "initiate"}
}

// initializes logger and miner
func initializeBlockchain() (int, Block) {
	fmt.Println("Program starting...")
	emptyBlock := Block{}
	genesisBlock := Genesis(&emptyBlock)
	fmt.Println("This code is reached.")
	minerLength := 1
	logChannel := make(chan Message, 1000000)
	mineChannels := make([]chan Block, minerLength)

	// initialize miner channels
	for i := 0; i < minerLength; i++ {
		mineChannels[i] = make(chan Block, 1)
	}

	// initialize logger
	logger = Logger{genesisBlock,
		sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)),
		&logChannel}

	// initialize miners
	miners = make([]Miner, minerLength)
	for i := 0; i < minerLength; i++ {
		miners[i] = Miner{string([]byte{byte(66 + i)}),
			genesisBlock,
			&mineChannels[i],
			sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)),
		}
	}
	return minerLength, genesisBlock
}

// initiate mining process
func startMining(minerLength int, genesisBlock Block) {
	start = time.Now()
	fmt.Println(genesisBlock.blockHeader.bits)
	fmt.Println(sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)))
	go logger.ListenForUpdate(miners)
	for i := 0; i < minerLength; i++ {
		go miners[i].Mine(logger)
	}
	time.Sleep(TimeLimit * time.Nanosecond)
}

func main() {
	minerLength, genesisBlock := initializeBlockchain()
	startMining(minerLength, genesisBlock)
}

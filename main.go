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

func PrettyPrintBlock(blocktoprint *Block) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println(" | version: \t\t", blocktoprint.blockHeader.version)
	fmt.Println(" | prevBlckHashPointer: ", blocktoprint.blockHeader.prevBlockHashPointer)
	fmt.Println(" | time:    \t\t", blocktoprint.blockHeader.time)
	fmt.Println(" | bits:    \t\t", blocktoprint.blockHeader.bits)
	fmt.Println(" | nonce:   \t\t", blocktoprint.blockHeader.nonce)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
}

func main() {
	fmt.Println("\nWelcome! Thank you for starting the program.\n")
	emptyBlock := Block{}
	genesisBlock := Genesis(&emptyBlock)
	fmt.Println("The first block has been created.")
	PrettyPrintBlock(&genesisBlock)
	minerLength := 3
	logChannel := make(chan Message, 1000000)
	mineChannels := make([]chan Block, minerLength)
	for i := 0; i < minerLength; i++ {
		mineChannels[i] = make(chan Block, 1)
	}
	logger = Logger{genesisBlock,
		sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)),
		&logChannel}
	miners = make([]Miner, minerLength)
	for i := 0; i < minerLength; i++ {
		miners[i] = Miner{string([]byte{byte(66 + i)}),
			genesisBlock,
			&mineChannels[i],
			sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)),
		}
	}
	start = time.Now()
	fmt.Println("\nMiners will now attempt to solve the puzzle given the following hash value:")
	fmt.Println(sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)), "\n")
	go logger.ListenForUpdate(miners)
	for i := 0; i < minerLength; i++ {
		go miners[i].Mine(logger)
	}
	time.Sleep(TimeLimit * time.Nanosecond)
}

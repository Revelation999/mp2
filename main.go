package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

var logger Logger
var miners []Miner

const TimeLimit = 300_000_000_000 // sets mining time limit to 5 minutes (300,000,000,000 nanoseconds)

var newpuzzle time.Time
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

func PrettyPrintBlock(blockToPrint *Block) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println(" | version: \t\t", blockToPrint.blockHeader.version)
	fmt.Println(" | prevBlockHashPointer: ", blockToPrint.blockHeader.prevBlockHashPointer)
	fmt.Println(" | time:    \t\t", blockToPrint.blockHeader.time)
	fmt.Println(" | bits:    \t\t", blockToPrint.blockHeader.bits)
	fmt.Println(" | nonce:   \t\t", blockToPrint.blockHeader.nonce)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
}

// initializes logger and miner
func initializeBlockchain() (int, Block) {
	fmt.Println("\nWelcome! Thank you for starting the program.\n")
	emptyBlock := Block{}
	genesisBlock := Genesis(&emptyBlock)
	fmt.Println("The first block has been created.")
	PrettyPrintBlock(&genesisBlock)
	fmt.Println("How many miners would you like to simulate? Please enter an integer value not equal to 0.")
    var minerLength int
    fmt.Scanln(&minerLength)
	for minerLength<=0{
		fmt.Println("Please enter a valid number of miners.")
		fmt.Scanln(&minerLength)
	}
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
	newpuzzle = time.Now()
	fmt.Println("\nMiners will now attempt to solve the puzzle given the following hash value:")
	fmt.Println(sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)), "\n")
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

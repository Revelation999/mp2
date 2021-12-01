package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var logger Logger
var miners []Miner

const TimeLimit = 300_000_000_000 // sets mining time limit to 5 minutes (300,000,000,000 nanoseconds)

var newpuzzle time.Time
var start time.Time

/*
	@input emptyBlock // An empty block that will serve as the first block of the chain
	@output Block // A block with an initial newHeader and the transaction "initiate"
	Genesis is a function that initializes the first block of a blockchain
*/
func Genesis(emptyBlock *Block) Block {
	var hash [32]byte
	for i := 0; i < 32; i++ {
		hash[i] = byte(rand.Intn(256))
	}
	fmt.Println("What would you like to set the difficulty level to? Enter n (0 < n < 32) in 2^(256-8(n)): ")
	var difficultyint int
	fmt.Scanln(&difficultyint)
	for difficultyint > 32 || difficultyint < 0 {
		fmt.Println("Please enter a valid n value for the difficulty (0 < n < 32).")
		fmt.Scanln(&difficultyint)
	}
	var difficulty [32]byte
	//for i := 0; i < 32; i++ {
	//	if i != difficultyint {
	//		difficulty[i] = 0
	//	} else {
	//		difficulty[i] = 1
	//	}
	//}
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

/*
	@input blockToPrint // A block that will be parsed to print
	PrettyPrintBlock is a function that will print out the contents of a block in a clean format
*/
func PrettyPrintBlock(blockToPrint *Block) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println(" | version: \t\t", blockToPrint.blockHeader.version)
	fmt.Println(" | prevBlockHashPointer: ", blockToPrint.blockHeader.prevBlockHashPointer)
	fmt.Println(" | time:    \t\t", blockToPrint.blockHeader.time)
	fmt.Println(" | bits:    \t\t", blockToPrint.blockHeader.bits)
	fmt.Println(" | nonce:   \t\t", blockToPrint.blockHeader.nonce)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println(" | transaction: \t", blockToPrint.transaction)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
}

/*
	@output minerLength // An integer value that will say how many miners there are
			genesisBlock // A Block instance that will serve as the first well constructed block in the blockchain
	InitializeBlockChain is a function that will ask the user how many miners they would like to simulate.
	This function will also initialize the logger, the miners, as well as the channel of the logger and those of the miners
*/
func initializeBlockchain() (int, Block) {
	fmt.Println("\nWelcome! Thank you for starting the program.\n")
	emptyBlock := Block{}
	genesisBlock := Genesis(&emptyBlock)
	fmt.Println("The first block has been created.")
	PrettyPrintBlock(&genesisBlock)
	fmt.Println("How many miners would you like to simulate? Please enter an integer value greater than or equal to 0.")
	var minerLength int
	fmt.Scanln(&minerLength)
	for minerLength <= 0 {
		fmt.Println("Please enter a valid number of miners.")
		fmt.Scanln(&minerLength)
	}
	logChannel := make(chan Message, 1000000)
	mineChannels := make([]chan Block, minerLength)

	// Initialize Miner Channels
	for i := 0; i < minerLength; i++ {
		mineChannels[i] = make(chan Block, 1)
	}

	// Initialize Logger object
	logger = Logger{genesisBlock,
		sha256.Sum256(HeaderToByteSlice(genesisBlock.blockHeader)),
		&logChannel}

	// Initialize miners
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

/*
	@input minerLength // The number of miners that the blockchain will simulate with
		   genesisBlock // A block that the Miners will derive their first puzzle from
	startMining is a function that will capture the current time after initiation and will initialize the goRoutines for both the miners and the logger
*/
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
	fmt.Printf("Max GMP value is %d", MaxParallelism())
	runtime.GOMAXPROCS() // Set GMP value for experiment
	minerLength, genesisBlock := initializeBlockchain()
	startMining(minerLength, genesisBlock)
}

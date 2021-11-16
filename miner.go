package main

import "fmt"

var globalNodeLog = make([]node, 5) // Keeps track of all nodes
var blockChain = make([]block, 5)   // Representation of the blockchain

type block struct {
	blockHeader blockHeader
	transaction string
}

type blockHeader struct {
	version              int
	hashPrevBlockHeader  string
	merkleRootHashFiller string
	time                 int
	bits                 int
	nonce                int
}

type node struct {
	branchLength int
	nodeChannel  chan int
}

// Initializes blocks
func initializeBlock(blockHeader blockHeader, transaction string) {
	newBlock := block{blockHeader, transaction}
	blockChain = append(blockChain, newBlock)
}

// Initializes block headers
func initializeBlockHeader(version int, hashPrevBlockHeader string,
	merkleRootHashFiller string, time int, bits int, nonce int) blockHeader {

	newBlockHeader := blockHeader{
		version, hashPrevBlockHeader,
		merkleRootHashFiller, time, bits, nonce}

	return newBlockHeader
}

// Initializes nodes
func initializeNode(nodeID int, branchLength int) {

	nodeChannel := make(chan int)
	newNode := node{branchLength, nodeChannel}
	globalNodeLog[nodeID] = newNode

}

// Primary initialize function that initializes blockHeaders, blocks, and nodes
func initialize() {
	fmt.Println("--------------------------------------------")
	fmt.Println("Welcome to Everything Distributed. How many nodes do you want")
	fmt.Println("If you would like to quit, please enter 'q'. ")
	nodeNumber := 0
	_, _ = fmt.Scanf("%d", nodeNumber)

	for counter := 0; counter < nodeNumber; counter++ {
		newBlockHeader := initializeBlockHeader(
			1, "prevHash",
			"merkleRoot", 5, 10, 2361345)
		initializeBlock(newBlockHeader, "Asher pays Lewis 5 bitcoin")
		initializeNode(counter, 0)
	}
}

func main() {
	initialize()
	fmt.Println("Finished intitializing")
	fmt.Println(globalNodeLog)
	fmt.Println(blockChain)

}

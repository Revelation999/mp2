# MP2
MP2 uses Go-channels and Go-routines to simulate mining and tamper-resistant log.
It implements Bitcoin: https://github.com/bitcoin/bitcoin

Authors: Steve Huang, Asher Kang, Maria Ringes. 

## How to Run 
### 1. Clone Github Repository

## Specification of Program Behavior

### Logger and Tamper-Resistant Log (Blockchain)
The Logger has 4 methods that define its behavior. 

1) `newBlock()` appends the newly approved block onto the blockchain. 
2) `UpdateBlock()` creates the next to-be-solved block on the blockchain and sends it to all the miners to be mined. 
3) `CheckNonce()` confirms whether the proposed puzzle solution fits the target difficulty value. 
4) `ListenForUpdate()` runs in a `Go-routine` and calls `CheckNonce()` on any proposed puzzle solution that is sent by a miner into the logger’s channel. 

### Miners and Mining (Puzzle Solving)
The Miners have 2 methods that define their behavior.

1) `Mine()` runs in a `Go-routine` and houses the main life-cycle of the miner. It calls `HasUpdate()` to check for new blocks from the logger. It repeatedly tries int values as a nonce to solve the puzzle, starting with value 1 and incrementing by 1. Finally, it terminates if the blockchain has been preserved for 5 minutes. 
2) `HasUpdate()` checks the miner’s channel for a new block sent from the logger. If it has received a new block, the method returns true. Otherwise, it returns false. 

### Supported Faulty Behavior 

#### Byzantine (Bogus solution)

Our program accounts for the Byzantine fault of a miner sending a bogus solution to the logger. The logger's CheckNonce() function returns a boolean false value of 0 if the miner's proposed solution does not solve the puzzle.

#### Crash Stop 

## Similarities to the Official Bitcoin Repository

### Block Headers

### Creation of Genesis Block

### Mining

## Screenshot
The following screenshot shows an example run with a set difficulty of 2^240. This example consists of 5 miners (B, C, D, E & F).
<img width="1113" alt="Screen Shot 2021-11-22 at 8 59 48 PM" src="https://user-images.githubusercontent.com/60116121/142960948-c31c652b-dfdb-4967-9712-397ee753a11a.png">
s 

## Workflow

## Custom Data Structures
```go

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
type Miner struct {
	identity      string
	block         Block
	mailbox       *chan Block
	currBlockHash [32]byte
}

type Logger struct {
	block         Block
	currBlockHash [32]byte
	mailbox       *chan Message
}

```

## Exit Codes 
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References 
Include notes about the crypto package we used 

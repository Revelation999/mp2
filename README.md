# MP2
MP2 uses Go-channels and Go-routines to simulate mining and tamper-resistant log.

Authors: Steve Huang, Asher Kang, Maria Ringes. 

## How to Run 
#### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/Revelation999/mp2`.

#### Step 2: Begin Bitcoin/Blockchain Implementation
Change the current directory into the recently cloned `mp2` folder. Start the Bitcoin/Blockchain protocol with `go run mp2`. 


#### Step 3: Interact with Command Line
A) **Difficulty Level** -- The program will ask the user to enter an `n` value between (and including) 0 and 32 such that the difficulty level is set to 2^(256-8n). The larger the value `n`, the smaller the difficulty level will be. A smaller difficulty level will make the puzzle harder for the miners to solve.

B) **Number of Miners** -- The program will ask the user how many miners to simulate in the blockchain. This integer should be greater than 0.

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

## Similarities to the Official Bitcoin Repository
The following are screenshots of code from the Bitcoin repository. These lines served as inspiration for our implementation of MP2. 

### Block Header Structure -- (From src>chain.h)
![blockheader](https://user-images.githubusercontent.com/15258611/142967019-4730c17c-e27d-4d4e-be98-49514f48f757.png)

### Creation of Genesis Block -- (From src>chainparams.cpp)
![genesis](https://user-images.githubusercontent.com/15258611/142967796-dab5b4a0-4121-429d-af29-9fc58e6ea7c1.png)

### SHA256 Hash Function -- (From src>hash.h)
![sha256](https://user-images.githubusercontent.com/15258611/142968387-d4863fdb-d5e6-42b7-85fd-42c7ed174e62.png)

### Mining: Nonce guess incrementation -- (From src>miner.cpp)
![Mining](https://user-images.githubusercontent.com/15258611/142966827-a2c33b27-936d-4319-8c36-e6970d1b9f74.png)

## Screenshot
The following screenshot shows an example run with a set difficulty of 2^240. This example consists of 5 miners (B, C, D, E & F).
<img width="1113" alt="Screen Shot 2021-11-22 at 8 59 48 PM" src="https://user-images.githubusercontent.com/60116121/142960948-c31c652b-dfdb-4967-9712-397ee753a11a.png">

The following screenshot shows an example run where the a miner proposes a nonce that does not satisfy the puzzle. In this example, we have reversed the compare statement such that as long as the hash value using the proposed nonce is greater than the difficulty, we send the guessed nonce value to the mailbox of the logger. As you can see, this does not force any block update and will just have the miner continue trying other values, hence the repetitive sends to the mailbox of the logger.
<img width="1113" alt="Screen Shot 2021-11-22 at 10 11 12 PM" src="https://user-images.githubusercontent.com/60116121/142965989-42721649-4112-4d80-9bbc-da681b8ef74d.png">

## Workflow
![MP2 - Main Workflow](https://user-images.githubusercontent.com/60116121/142963455-08cd1f29-0789-4d64-9d5b-04bc609093ea.png)


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
crypto/sha256

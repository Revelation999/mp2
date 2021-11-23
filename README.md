# MP2
MP2 uses Go-channels and Go-routines to simulate mining and tamper-resistant log.
It implements Bitcoin: https://github.com/bitcoin/bitcoin

Authors: Steve Huang, Asher Kang, Maria Ringes. 

## How to Run 
### 1. Clone Github Repository

## Specification of Program Behavior

### Logger

### Miners 

### Puzzle Solving (Mining)

### Tamper-resistant Log (Blockchain)

## Screenshots 

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

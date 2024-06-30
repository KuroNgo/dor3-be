package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// calculateHash Convert the block's data to JSON
// Concatenated the previous block's hash, and the current block's data, timestamp, and Pow
// Hashed the earlier concatenation with the SHA256 algorithm
// Return the base 16 hash as a string
func (b *Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String()
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash) // %x is hexadecimals
}

// mine for our Block type that repeatedly increments the Pow value and calculates the block hash
// until we get a valid one
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// CreateBlockchain create a genesis block ( the first block on the blockchain) for our blockchain and
// returns a new instance of the Blockchain type. Add the following code to the blockchain.go file
// We set the hash of our genesis block to 0. Because it is the first block in the blockchain,
// there is no value for the previous hash, and the data property is empty. Then we create a new instance of the Blockchain type
// and stored the genesis block along with the blockchain's difficulty
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

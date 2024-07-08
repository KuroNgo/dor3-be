package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Timestamp     int64  `bson:"timestamp" json:"timestamp"`
	Data          int64  `bson:"data" json:"data"`
	PrevBlockHash []byte `bson:"prevBlockHash" json:"prevBlockHash"`
	Hash          []byte `bson:"hash" json:"hash"`
}

// BlockChain represents the blockchain
type BlockChain struct {
	Blocks []*Block
}

// SetHash calculates and sets the block's hash
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, []byte(strconv.FormatInt(b.Data, 10)), timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// NewBlock creates and returns a Block
func NewBlock(data int64, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), data, prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// AddBlock saves the block into the blockchain
func (bc *BlockChain) AddBlock(data int64) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// NewGenesisBlock creates and returns the genesis block
func NewGenesisBlock() *Block {
	return NewBlock(0, []byte{})
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

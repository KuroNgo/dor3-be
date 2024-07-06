package blockchain

//
//type Block struct {
//	data         map[string]interface{}
//	hash         string
//	previousHash string
//	timestamp    time.Time
//	pow          int
//}
//
//type Blockchain struct {
//	genesisBlock Block
//	chain        []Block
//	difficulty   int
//}
//
//// CalculateHash Convert the block's data to JSON
//// Concatenated the previous block's hash, and the current block's data, timestamp, and Pow
//// Hashed the earlier concatenation with the SHA256 algorithm
//// Return the base 16 hash as a string
//func (b *Block) CalculateHash() string {
//	data, _ := json.Marshal(b.data)
//	blockData := b.previousHash + string(data) + b.timestamp.String()
//	blockHash := sha256.Sum256([]byte(blockData))
//	return fmt.Sprintf("%x", blockHash) // %x is hexadecimals
//}
//
//// Mine for our Block type that repeatedly increments the Pow value and calculates the block hash
//// until we get a valid one
//func (b *Block) Mine(difficulty int) {
//	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
//		b.pow++
//		b.hash = b.CalculateHash()
//	}
//}
//
//// CreateBlockchain create a genesis block ( the first block on the blockchain) for our blockchain and
//// returns a new instance of the Blockchain type. Add the following code to the blockchain.go file
//// We set the hash of our genesis block to 0. Because it is the first block in the blockchain,
//// there is no value for the previous hash, and the data property is empty. Then we create a new instance of the Blockchain type
//// and stored the genesis block along with the blockchain's difficulty
//func CreateBlockchain(difficulty int) Blockchain {
//	genesisBlock := Block{
//		hash:      "0",
//		timestamp: time.Now(),
//	}
//	return Blockchain{
//		genesisBlock,
//		[]Block{genesisBlock},
//		difficulty,
//	}
//}
//
//// AddBlockForTransaction collects the details of a transaction (sender, receiver, and transfer amount
//// Creates a new block with the transaction details
//// Mines the new block with the previous block hash, current block data, and generated Pow
//// Adds the newly created block to the blockchain
//func (b *Blockchain) AddBlockForTransaction(from, to string, amount float64) {
//	blockData := map[string]interface{}{
//		"from":   from,
//		"to":     to,
//		"amount": amount,
//	}
//
//	lastBlock := b.chain[len(b.chain)-1]
//	newBlock := Block{
//		data:         blockData,
//		previousHash: lastBlock.hash,
//		timestamp:    time.Now(),
//	}
//	newBlock.Mine(b.difficulty)
//	b.chain = append(b.chain, newBlock)
//}
//
//// IsValid the method is recalculated the hash of every block, compared them with the stored hash values
//// of the other blocks, and checked if the previous hash value of any other block is equal to
//// the hash value of the block before it. If any of the checks fail, the blockchain has been tampered with
//func (b *Blockchain) IsValid() bool {
//	for i := range b.chain[1:] {
//		previousBlock := b.chain[i]
//		currentBlock := b.chain[i+1]
//		if currentBlock.hash != currentBlock.CalculateHash() || currentBlock.previousHash != previousBlock.hash {
//			return false
//		}
//	}
//
//	return true
//}
//
//// # Adding new nodes to the blockchain
////
//// A blockchain isn't useful to anyone if it's confined to just one device. So, you might be
//// wondering how to add new nodes to the blockchain prototype you just built.
//// Pulling this off involves creating a Peer to Peer P2P network to connect the nodes, implementing a consensus algorithm, and taking security countermeasures.

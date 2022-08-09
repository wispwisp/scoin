package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/wispwisp/scoin/transaction"
)

type Block struct {
	Index        int
	PrevHash     string
	Nonce        int
	Transactions []transaction.Transaction
}

func (block *Block) String() string {
	j, err := json.Marshal(*block)
	if err != nil {
		panic("Can't convert Block to json")
	}
	return string(j)
}

func (block *Block) Hash() string {
	j, err := json.Marshal(*block)
	if err != nil {
		panic("Can't convert Block to json")
	}

	h := sha256.New()
	h.Write([]byte(string(j)))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateNextBlock(lastBlock *Block, transactions []transaction.Transaction) Block {
	var nextBlock Block

	nextBlock.Index = lastBlock.Index + 1
	nextBlock.Nonce = lastBlock.Nonce
	nextBlock.PrevHash = lastBlock.Hash()
	nextBlock.Transactions = transactions

	return nextBlock
}

func (block *Block) GenerateNextBlock(transactions []transaction.Transaction) Block {
	return GenerateNextBlock(block, transactions)
}

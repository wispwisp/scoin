package mine

import (
	"log"
	"math/rand"
	"time"

	"github.com/wispwisp/scoin/block"
	"github.com/wispwisp/scoin/transaction"
)

const maxAccumulateTransactions = 5
const waitToProcessTransactionsAnywayInSeconds = 10

// Accumulate enough transactions (or process it anyway if no transaction for a period of time)
func accumulateTransactions(
	transactionsChan chan transaction.Transaction,
) (transactions []transaction.Transaction) {
	counter := 0
	for {
		select {
		case transaction := <-transactionsChan:
			transactions = append(transactions, transaction)
			counter++
			if counter > maxAccumulateTransactions {
				return
			}
		case <-time.After(waitToProcessTransactionsAnywayInSeconds * time.Second):
			if len(transactions) > 0 {
				return
			}
		}
	}
}

func ProofOfWork(block *block.Block, nonce int) (int, bool) {
	block.Nonce = nonce

	h := block.Hash()
	found := (h[0] == '0') && (h[1] == '0')

	return nonce, found
}

func Mine(
	blockchain *[]block.Block,
	transactionsChan chan transaction.Transaction,
	consensusChan chan block.Block,
) {
	log.Println("Mine")

	for {
		transactions := accumulateTransactions(transactionsChan)

		currentBlockchainLen := len(*blockchain)
		lastBlock := (*blockchain)[currentBlockchainLen-1]
		nextBlock := lastBlock.GenerateNextBlock(transactions)

		// Find nonce
		nonce := rand.Int()
		found := false
		for {
			// Check if other node found consensus earlier
			if currentBlockchainLen != len(*blockchain) {
				log.Println("currentBlockchainLen != len(*blockchain): ", currentBlockchainLen, len(*blockchain))
				break
			}

			nonce, found = ProofOfWork(&nextBlock, nonce+1)
			if found {
				log.Println("Found nonce, send to blockchain: ", nonce)
				consensusChan <- nextBlock
				break
			}
		}
	}
}
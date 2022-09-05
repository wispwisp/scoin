package mine

import (
	"log"
	"math/rand"
	"time"

	"github.com/wispwisp/scoin/block"
	"github.com/wispwisp/scoin/nethelpers"
	"github.com/wispwisp/scoin/node"
	"github.com/wispwisp/scoin/transaction"
)

const maxAccumulateTransactions = 5
const waitToProcessTransactionsAnywayInSeconds = 10

func isUnique(transactions *[]transaction.Transaction, t *transaction.Transaction) bool {
	for _, s := range *transactions {
		if s == *t {
			return false
		}
	}
	return true
}

// Accumulate enough transactions (or process it anyway if no transaction for a period of time)
func accumulateTransactions(
	nodesInfo *node.NodesInfo,
	transactionsChan chan transaction.Transaction,
) (transactions []transaction.Transaction) {
	counter := 0
	for {
		select {
		case transaction := <-transactionsChan:
			if isUnique(&transactions, &transaction) {
				transactions = append(transactions, transaction)
				nethelpers.SendTransactionToOtherNodes(nodesInfo, &transaction)
			}
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
	blockchain *block.Blockchain,
	nodesInfo *node.NodesInfo,
	transactionsChan chan transaction.Transaction,
	consensusChan chan block.Block,
) {
	log.Println("Mine")

	for {
		transactions := accumulateTransactions(nodesInfo, transactionsChan)

		currentBlockchainLen := blockchain.Len()
		lastBlock := blockchain.GetLastBlock()
		nextBlock := lastBlock.GenerateNextBlock(transactions)

		// Find nonce
		nonce := rand.Int()
		found := false
		for {
			// Check if other node found consensus earlier
			if currentBlockchainLen != blockchain.Len() {
				log.Println("Other node found consesnus, drop current mining")
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

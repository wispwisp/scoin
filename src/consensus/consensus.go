package consensus

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wispwisp/scoin/block"
	"github.com/wispwisp/scoin/node"
)

// TODO:
// 1) Ask all nodes in network for their blockchains
// 2) Validate new blocks if blockchain differ
// 3) Exclude mined transactions from pool
// 4) Update own blockchain

// TODO:
// if we have consensus, node should drop transactions
// which has been added to a new block.

func requestForNode(uri string) (blockchainPart []block.Block, success bool) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("Failed to request node ", uri, ", error ", err)
		return
	}
	defer resp.Body.Close()

	// TODO If resp not succsess: log and exit !

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error parsing response body:", err)
		return
	}

	err = json.Unmarshal(body, &blockchainPart)
	if err != nil {
		log.Println("error parsing blockchain part:", err)
		return
	}

	success = true
	return
}

func checkNode(uri string, blockchain *block.Blockchain) (needUpdate bool, blockchainPart []block.Block) {
	log.Println("Check for node: ", uri)

	blockchainPart, success := requestForNode(uri)
	if !success {
		return
	}

	lenOfRest := len(blockchainPart)
	if lenOfRest == 0 {
		log.Println("recieved empty blockchain from ", uri)
		return
	}

	block := blockchainPart[0]
	lastBlock := blockchain.GetLastBlock()

	// Validate previous blockchain by checking hashes
	if block.PrevHash != lastBlock.PrevHash {
		log.Println("Incorrect blockchain for node: ", uri)
		return
	}

	if lenOfRest != 1 {
		needUpdate = true
		blockchainPart = blockchainPart[1:] // Remove starting block of blockchain part
	}

	return
}

func getLongestBlockchainIndex(blockchains *[][]block.Block) (maxLenght int, maxBlockChainIndex int) {
	for i, bc := range *blockchains {
		l := len(bc)
		if l > maxLenght {
			maxLenght = l
			maxBlockChainIndex = i
		}
	}
	return
}

func consensusIteration(blockchain *block.Blockchain, nodesInfo *node.NodesInfo, consensusChan chan block.Block) {
	log.Println("Check other nodes...")

	// 1) Ask all nodes in network for their blockchains
	index := blockchain.Len() - 1

	var blockchains [][]block.Block
	for _, nodeInfo := range nodesInfo.Get() {
		uri := "http://" + nodeInfo.Uri + "/blockchain/" + strconv.Itoa(index)
		needUpdate, blockchainPart := checkNode(uri, blockchain)
		if needUpdate {
			blockchains = append(blockchains, blockchainPart)
		}
	}

	// Any blockchains to update current node?
	recievedFromOtherNode := false
	if len(blockchains) != 0 {
		_, maxBlockChainIndex := getLongestBlockchainIndex(&blockchains)
		// TODO: validate recieved blockchain
		blockchain.AddBlocks(blockchains[maxBlockChainIndex]...)
		recievedFromOtherNode = true // drop current mined block futher (if any)
		log.Println("Current blockchain updated from other node. Update blocks:", blockchains[maxBlockChainIndex])
	}

	// Check for current node POW result
	select {
	case nextBlock := <-consensusChan:
		if !recievedFromOtherNode {
			log.Println("Block from current node POW recieved, update blockchain")
			blockchain.Add(&nextBlock)
		} else {
			log.Println("Current blockchain updated from other node, drop current mined block")
		}
	default: // No block from current node
	}
}

func Consensus(blockchain *block.Blockchain, nodesInfo *node.NodesInfo, consensusChan chan block.Block) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		<-ticker.C
		consensusIteration(blockchain, nodesInfo, consensusChan)
	}
}

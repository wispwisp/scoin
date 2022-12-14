package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/wispwisp/scoin/block"
	"github.com/wispwisp/scoin/consensus"
	"github.com/wispwisp/scoin/mine"
	"github.com/wispwisp/scoin/nethelpers"
	"github.com/wispwisp/scoin/node"
	"github.com/wispwisp/scoin/transaction"
)

type Args struct {
	Port           *string
	InitBlockchain *bool
	InitalNodeAddr *string
}

func registerArgs() (args Args) {
	args.Port = flag.String("port", "8090", "server port")
	args.InitBlockchain = flag.Bool("init", false, "make initial trasaction")
	args.InitalNodeAddr = flag.String("node", "", "node addr for initialization")
	flag.Parse()
	return
}

func DefaultInitBlockchain(blockchain *block.Blockchain) {
	log.Println("Make initial transaction, create first block")
	firstTransaction := transaction.Transaction{From: "network", To: "addr1", Amount: 50}
	blockchain.Add(&block.Block{
		Index: 0, PrevHash: "none", Nonce: 0,
		Transactions: []transaction.Transaction{firstTransaction}})
}

func SetBlockchainFromOtherNode(blockchain *block.Blockchain, nodesInfo *node.NodesInfo) bool {
	log.Println("Sync blockchain from other nodes...")
	otherBlockchain := nethelpers.GetLongestBlockchainFromNodes(nodesInfo)
	if len(otherBlockchain) == 0 {
		log.Println("Blockchains not found.")
		return false
	}

	blockchain.AddBlocks(otherBlockchain...)
	return true
}

func main() {
	args := registerArgs()

	fileName := "../conf/nodes.json"
	var nodesInfo node.NodesInfo
	if err := nodesInfo.LoadFromFile(fileName); err != nil {
		log.Println("Error loading from", fileName, "error:", err)
	}

	var blockchain block.Blockchain

	if len(*args.InitalNodeAddr) != 0 {
		nodesInfo.Add(&node.NodeInfo{Uri: *args.InitalNodeAddr})
		log.Println("Add Nodes info with URL:", *args.InitalNodeAddr)
	}

	if *args.InitBlockchain {
		DefaultInitBlockchain(&blockchain)
	} else {
		if !SetBlockchainFromOtherNode(&blockchain, &nodesInfo) {
			return
		}
	}

	// Start mining
	transactionsChan := make(chan transaction.Transaction, 100)
	consensusChan := make(chan block.Block)
	go mine.Mine(&blockchain, &nodesInfo, transactionsChan, consensusChan)
	go consensus.Consensus(&blockchain, &nodesInfo, consensusChan)

	// http handlers:

	http.HandleFunc("/transaction", func(w http.ResponseWriter, req *http.Request) {
		log.Println("'/transaction' HTTP handler - add transaction.")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("error parsing request body:", err)
			http.NotFound(w, req)
			return
		}

		// var jsonRes map[string]interface{}
		// err = json.Unmarshal(body, &jsonRes)

		var t transaction.Transaction
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error parsing request's json transaction:", err)
			http.NotFound(w, req)
			return
		}

		transactionsChan <- t
	})

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, req *http.Request) {
		log.Println("'/blockchain' HTTP handler - show blockchain to clinet.")
		if encodeErr := json.NewEncoder(w).Encode(blockchain.Get()); encodeErr != nil {
			log.Println("Encode blockchain to json failed, err: ", encodeErr)
			http.NotFound(w, req)
			return
		}
	})

	// TODO: routers
	http.HandleFunc("/blockchain/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("'" + req.URL.Path + "' HTTP handler - show blockchain to clinet.")

		// Extract blockchain index from API
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect '/blockchain/<index>", http.StatusBadRequest)
			return
		}

		index, err := strconv.Atoi(pathParts[1])
		if err != nil {
			log.Println("Invalid block index requested: ", index, " error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if index >= blockchain.Len() {
			log.Println("Invalid block index requested: ", index)
			http.Error(w, "Invalid block index requested", http.StatusBadRequest)
			return
		}

		blockchainPart := blockchain.GetSlice(index)

		if encodeErr := json.NewEncoder(w).Encode(blockchainPart); encodeErr != nil {
			log.Println("Encode blockchain to json failed, err: ", encodeErr)
			http.NotFound(w, req)
			return
		}
	})

	http.HandleFunc("/addnode", func(w http.ResponseWriter, req *http.Request) {
		log.Println("'/addnode' HTTP handler - add addnode")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var nodeInfo node.NodeInfo
		err = json.Unmarshal(body, &nodeInfo)
		if err != nil {
			log.Println("error parsing node info:", err)
			http.Error(w, "error parsing node info", http.StatusBadRequest)
			return
		}

		nodesInfo.Add(&nodeInfo)
	})

	log.Println("Server started on", *args.Port, "port")
	http.ListenAndServe(":"+*args.Port, nil)
}

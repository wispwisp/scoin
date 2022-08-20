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
	"github.com/wispwisp/scoin/node"
	"github.com/wispwisp/scoin/transaction"
)

type Args struct {
	Port           *string
	InitBlockchain *bool
}

func registerArgs() (args Args) {
	args.Port = flag.String("port", "8090", "server port")
	args.InitBlockchain = flag.Bool("init", false, "make initial trasaction")
	flag.Parse()
	return
}

func main() {
	var blockchain []block.Block

	args := registerArgs()

	fileName := "../conf/nodes.json"
	var nodesInfo node.NodesInfo
	if err := nodesInfo.LoadFromFile(fileName); err != nil {
		log.Println("Error loading from ", fileName, ": ", err)
		return
	}

	if *args.InitBlockchain {
		log.Println("Make initial transaction, create first block")
		firstTransaction := transaction.Transaction{From: "network", To: "addr1", Amount: 50}
		b := block.Block{Index: 0, PrevHash: "none", Nonce: 0, Transactions: []transaction.Transaction{firstTransaction}}
		blockchain = append(blockchain, b)
	} else {
		// TODO
		log.Println("Sync blockchain from other nodes")
		log.Println("UNIMPLEMENTED. extiting")
		return
	}

	// Start mining
	transactionsChan := make(chan transaction.Transaction, 100)
	consensusChan := make(chan block.Block)
	go mine.Mine(&blockchain, transactionsChan, consensusChan)
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
		// TODO: send that transaction for all other nodes in blockchain
	})

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, req *http.Request) {
		log.Println("'/blockchain' HTTP handler - show blockchain to clinet.")
		if encodeErr := json.NewEncoder(w).Encode(blockchain); encodeErr != nil {
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

		if index >= len(blockchain) {
			log.Println("Invalid block index requested: ", index)
			http.Error(w, "Invalid block index requested", http.StatusBadRequest)
			return
		}

		blockchainPart := blockchain[index:]

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

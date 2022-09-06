package nethelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/wispwisp/scoin/block"
	"github.com/wispwisp/scoin/node"
	"github.com/wispwisp/scoin/transaction"
)

func SendTransactionToOtherNodes(nodesInfo *node.NodesInfo, t *transaction.Transaction) {
	ni := nodesInfo.Get()
	for _, nodeInfo := range ni {
		go func(uri string) {
			log.Println("Send transaction to", uri)

			jsonData, err := json.Marshal(t)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Println(err)
				return
			}

			// Log result
			// var res map[string]interface{}
			// json.NewDecoder(resp.Body).Decode(&res)
			// log.Println(res["json"])
		}("http://" + nodeInfo.Uri + "/transaction")
	}
}

func RequestForNode(uri string) (blockchainPart []block.Block, success bool) {
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

func GetLongestBlockchainIndex(blockchains *[][]block.Block) (maxLenght int, maxBlockchainIndex int) {
	for i := 0; i < len(*blockchains); i++ {
		l := len((*blockchains)[i])
		if l > maxLenght {
			maxLenght = l
			maxBlockchainIndex = i
		}
	}
	return
}

func GetLongestBlockchainFromNodes(nodesInfo *node.NodesInfo) []block.Block {
	ni := nodesInfo.Get()
	sz := len(ni)
	if sz == 0 {
		return nil
	}

	blockchains := make([][]block.Block, sz)

	var wg sync.WaitGroup

	for i, nodeInfo := range ni {
		wg.Add(1)
		go func(uri string, i int) {
			defer wg.Done()
			blockchainPart, success := RequestForNode(uri)
			if success {
				blockchains[i] = blockchainPart
			}
		}("http://"+nodeInfo.Uri+"/blockchain", i)
	}

	wg.Wait()

	_, maxIndex := GetLongestBlockchainIndex(&blockchains)

	return blockchains[maxIndex]
}

package nethelpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

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

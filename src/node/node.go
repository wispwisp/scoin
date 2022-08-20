package node

import (
	"sync"
	"os"
	"encoding/json"
)

type NodeInfo struct {
	Uri string `json:"uri"`
}

type NodesInfo struct {
	mtx       sync.Mutex
	nodesInfo []NodeInfo
}

func (ni *NodesInfo) Add(nodeInfo *NodeInfo) {
	ni.mtx.Lock()
	defer ni.mtx.Unlock()
	ni.nodesInfo = append(ni.nodesInfo, *nodeInfo)
}

func (ni *NodesInfo) Get() []NodeInfo { // Copy
	ni.mtx.Lock()
	defer ni.mtx.Unlock()
	return ni.nodesInfo
}


func (ni *NodesInfo) LoadFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	ni.mtx.Lock()
	defer ni.mtx.Unlock()

	err = json.Unmarshal(data, &ni.nodesInfo)
	if err != nil {
		return err
	}

	return nil
}

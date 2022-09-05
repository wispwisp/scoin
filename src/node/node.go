package node

import (
	"encoding/json"
	"os"
	"sync"
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

func (ni *NodesInfo) Get() []NodeInfo {
	ni.mtx.Lock()
	defer ni.mtx.Unlock()

	aCopy := make([]NodeInfo, len(ni.nodesInfo))
	copy(aCopy, ni.nodesInfo)

	return aCopy
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

package block

import "sync"

type Blockchain struct {
	mtx    sync.Mutex
	blocks []Block
}

func (bc *Blockchain) Len() int {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()
	return len(bc.blocks)
}

func (bc *Blockchain) Add(block *Block) {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()
	bc.blocks = append(bc.blocks, *block)
}

func (bc *Blockchain) AddBlocks(blocks ...Block) {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()
	bc.blocks = append(bc.blocks, blocks...)
}

func (bc *Blockchain) Get() (aCopy []Block) {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	aCopy = make([]Block, len(bc.blocks))
	copy(aCopy, bc.blocks)

	return
}

func (bc *Blockchain) GetSlice(from int) (aCopy []Block) {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	if from >= len(bc.blocks) {
		return
	}

	blockchainPart := bc.blocks[from:]

	aCopy = make([]Block, len(blockchainPart))
	copy(aCopy, blockchainPart)

	return
}

func (bc *Blockchain) GetLastBlock() Block {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) Validate() bool {
	bc.mtx.Lock()
	aCopy := make([]Block, len(bc.blocks))
	copy(aCopy, bc.blocks)
	bc.mtx.Unlock()

	// TODO: Iterate over block, make hash and validate prevHash

	return true
}

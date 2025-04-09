package main

type Block struct {
	EntityList []EntityInterface
}

func NewBlock() *Block {
	return &Block{
		EntityList: make([]EntityInterface, 0),
	}
}

func (b *Block) AddEntity(e EntityInterface) {
	if b.EntityList == nil {
		b.EntityList = make([]EntityInterface, 0)
	}
	b.EntityList = append(b.EntityList, e)
}

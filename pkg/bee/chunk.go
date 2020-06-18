package bee

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	maxChunkSize = 4096
	spanInfoSize = 8
)

// Chunk represents Bee chunk
type Chunk struct {
	address swarm.Address
	data    []byte
	span    int
}

// NewChunk returns new chunk
func NewChunk(data []byte) (Chunk, error) {
	if len(data) > maxChunkSize {
		return Chunk{}, fmt.Errorf("create chunk: requested size too big (max %d bytes)", maxChunkSize)
	}

	return Chunk{data: data}, nil
}

// NewRandomChunk returns new pseudorandom chunk
func NewRandomChunk(r *rand.Rand) (c Chunk, err error) {
	data := make([]byte, r.Intn(maxChunkSize-spanInfoSize))
	if _, err := r.Read(data); err != nil {
		return Chunk{}, fmt.Errorf("create random chunk: %w", err)
	}

	span := len(data)
	b := make([]byte, spanInfoSize)
	binary.LittleEndian.PutUint64(b, uint64(span))
	data = append(b, data...)

	c = Chunk{data: data, span: span}
	return
}

// Address returns chunk's address
func (c *Chunk) Address() swarm.Address {
	return c.address
}

// Data returns chunk's data
func (c *Chunk) Data() []byte {
	return c.data
}

// Size returns chunk size
func (c *Chunk) Size() int {
	return len(c.data)
}

// Span returns chunk span
func (c *Chunk) Span() int {
	return c.span
}

// ClosestNode returns chunk's closest node of a given set of nodes
func (c *Chunk) ClosestNode(nodes []swarm.Address) (closest swarm.Address, err error) {
	closest = nodes[0]
	for _, a := range nodes[1:] {
		dcmp, err := swarm.DistanceCmp(c.Address().Bytes(), closest.Bytes(), a.Bytes())
		if err != nil {
			return swarm.Address{}, fmt.Errorf("find closest node: %w", err)
		}
		switch dcmp {
		case 0:
			// do nothing
		case -1:
			// current node is closer
			closest = a
		case 1:
			// closest is already closer to chunk
			// do nothing
		}
	}

	return
}

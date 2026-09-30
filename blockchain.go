package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash"
	"strings"
	"time"
)

const maxDifficulty = 6

var ErrInvalidDifficulty = errors.New("difficulty must be between 1 and 6")

type Block struct {
	Index        uint64
	Timestamp    int64
	Data         string
	PreviousHash string
	Hash         string
	Nonce        uint64
}

type Chain struct {
	blocks     []Block
	difficulty int
}

func NewChain(difficulty int) (*Chain, error) {
	if difficulty < 1 || difficulty > maxDifficulty {
		return nil, ErrInvalidDifficulty
	}

	chain := &Chain{difficulty: difficulty}
	chain.blocks = append(chain.blocks, mineBlock(0, "Genesis Block", "", difficulty))
	return chain, nil
}

func (c *Chain) AddBlock(data string) Block {
	previous := c.blocks[len(c.blocks)-1]
	block := mineBlock(previous.Index+1, data, previous.Hash, c.difficulty)
	c.blocks = append(c.blocks, block)
	return block
}

func (c *Chain) Blocks() []Block {
	return append([]Block(nil), c.blocks...)
}

func (c *Chain) Validate() bool {
	target := strings.Repeat("0", c.difficulty)
	for i, block := range c.blocks {
		expectedPrevious := ""
		if i > 0 {
			expectedPrevious = c.blocks[i-1].Hash
			if block.Index != c.blocks[i-1].Index+1 {
				return false
			}
		} else if block.Index != 0 {
			return false
		}

		if block.PreviousHash != expectedPrevious ||
			block.Hash != calculateHash(block) ||
			block.Hash[:c.difficulty] != target {
			return false
		}
	}
	return len(c.blocks) > 0
}

func mineBlock(index uint64, data, previousHash string, difficulty int) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().UnixNano(),
		Data:         data,
		PreviousHash: previousHash,
	}
	target := strings.Repeat("0", difficulty)
	for {
		block.Hash = calculateHash(block)
		if block.Hash[:difficulty] == target {
			return block
		}
		block.Nonce++
	}
}

func calculateHash(block Block) string {
	h := sha256.New()
	writeUint64(h, block.Index)
	writeUint64(h, uint64(block.Timestamp))
	writeString(h, block.Data)
	writeString(h, block.PreviousHash)
	writeUint64(h, block.Nonce)
	return hex.EncodeToString(h.Sum(nil))
}

func writeUint64(h hash.Hash, value uint64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], value)
	_, _ = h.Write(encoded[:])
}

func writeString(h hash.Hash, value string) {
	writeUint64(h, uint64(len(value)))
	_, _ = h.Write([]byte(value))
}

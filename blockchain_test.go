package blockchain

import (
	"errors"
	"testing"
)

func TestNewChainRejectsInvalidDifficulty(t *testing.T) {
	for _, difficulty := range []int{0, maxDifficulty + 1} {
		if _, err := NewChain(difficulty); !errors.Is(err, ErrInvalidDifficulty) {
			t.Fatalf("NewChain(%d) error = %v, want %v", difficulty, err, ErrInvalidDifficulty)
		}
	}
}

func TestAddBlockMinesAndLinksBlocks(t *testing.T) {
	chain, err := NewChain(2)
	if err != nil {
		t.Fatal(err)
	}

	block := chain.AddBlock("rental created")
	blocks := chain.Blocks()
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if block.Index != 1 || block.PreviousHash != blocks[0].Hash {
		t.Fatalf("new block is not linked to genesis: %+v", block)
	}
	if block.Hash[:2] != "00" {
		t.Fatalf("mined hash %q does not satisfy difficulty", block.Hash)
	}
	if !chain.Validate() {
		t.Fatal("valid chain failed validation")
	}
}

func TestValidateDetectsTampering(t *testing.T) {
	chain, err := NewChain(1)
	if err != nil {
		t.Fatal(err)
	}
	chain.AddBlock("rental created")
	chain.blocks[1].Data = "rental altered"

	if chain.Validate() {
		t.Fatal("tampered chain passed validation")
	}
}

func TestBlocksReturnsCopy(t *testing.T) {
	chain, err := NewChain(1)
	if err != nil {
		t.Fatal(err)
	}
	blocks := chain.Blocks()
	blocks[0].Data = "altered"

	if chain.Blocks()[0].Data == "altered" {
		t.Fatal("modifying returned blocks changed the chain")
	}
}

# Go Blockchain Study Project

A small, standard-library-only proof-of-work blockchain for learning the core
ideas behind blocks, hashing, mining, and chain validation.

## Use

```go
chain, err := blockchain.NewChain(3)
if err != nil {
	log.Fatal(err)
}

block := chain.AddBlock("rental created")
fmt.Println(block.Hash)
fmt.Println(chain.Validate())
```

Import the package with `github.com/bergsantana/car-rental-blockchain`. The
difficulty is the number of leading zeroes required in each block's SHA-256
hash; values from 1 to 6 are accepted. `Blocks` returns a copy of the chain's
blocks.

Run the tests with `go test ./...`.
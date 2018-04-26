/*
 * @Author: willxm
 * @Date: 2018-04-27 00:03:51
 * @Last Modified by: willxm
 * @Last Modified time: 2018-04-27 00:06:18
 */
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block is BlockChain item
type Block struct {
	Index     int
	Timestamp int64
	Data      string
	Hash      string
	PrevHash  string
}

// Blockchain is chain
var Blockchain []Block

// Calculate the block's hash
func CalculateHash(block Block) string {
	record := string(block.Index) + string(block.Timestamp) + string(block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Generate a new block by previous block
func GenerateBlock(oldBlock Block, data string) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.Unix()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}

// Check block is valid
func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// keep the new blockchain is the longest
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

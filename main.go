/*
 * @Author: willxm
 * @Date: 2018-04-26 00:24:13
 * @Last Modified by: ???@willxmwillxm
 * @Last Modified time: 2018-04-26 01:07:08
 */
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

// Message is the http request
type Message struct {
	Data string
}

// Blockchain is chain
var Blockchain []Block

// Calculate the block's hash
func calculateHash(block Block) string {
	record := string(block.Index) + string(block.Timestamp) + string(block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Generate a new block by previous block
func generateBlock(oldBlock Block, data string) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.Unix()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// Check block is valid
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// keep the new blockchain is the longest
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.Data)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}
func getBlock(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	//init creation block
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.Unix(), "", "", ""}
		Blockchain = append(Blockchain, genesisBlock)
	}()

	http.HandleFunc("/write", writeBlock)
	http.HandleFunc("/get", getBlock)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

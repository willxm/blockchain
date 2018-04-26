/*
 * @Author: willxm
 * @Date: 2018-04-26 00:24:13
 * @Last Modified by: willxm
 * @Last Modified time: 2018-04-27 00:08:07
 */
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/willxm/blockchain/core"
)

// Message is the http request
type Message struct {
	Data string
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := core.GenerateBlock(core.Blockchain[len(core.Blockchain)-1], m.Data)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if core.IsBlockValid(newBlock, core.Blockchain[len(core.Blockchain)-1]) {
		newBlockchain := append(core.Blockchain, newBlock)
		core.ReplaceChain(newBlockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}
func getBlock(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(core.Blockchain, "", "  ")
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
		genesisBlock := core.Block{0, t.Unix(), "", "", ""}
		core.Blockchain = append(core.Blockchain, genesisBlock)
	}()

	http.HandleFunc("/write", writeBlock)
	http.HandleFunc("/get", getBlock)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

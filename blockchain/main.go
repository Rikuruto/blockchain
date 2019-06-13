package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var bc blockchain

type blockchain struct {
	chain []block
}

type block struct {
	Index        int    `json:"index,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	Proof        int    `json:"proof,omitempty"`
	PreviousHash string `json:"previous_hash,omitempty"`
}

func (bc *blockchain) CreateBlock(proof int, previousHash string) block {
	blk := block{
		Index:        len(bc.chain) + 1,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05.000000"),
		Proof:        proof,
		PreviousHash: previousHash,
	}
	bc.chain = append(bc.chain, blk)
	return blk
}

func (bc blockchain) GetPreviousBlock() block {
	return bc.chain[len(bc.chain)-1]
}

func (bc blockchain) proofOfWork(previousProof int) int {
	newProof := 1
	check := false
	for !check {
		h := sha256.New()
		h.Write([]byte(strconv.Itoa(newProof*newProof - previousProof*previousProof)))
		bs := h.Sum(nil)
		if strings.HasPrefix(fmt.Sprintf("%x", bs), "0000") {
			check = true
		} else {
			newProof++
		}
	}
	return newProof
}

func (blk block) Hash() string {
	j, _ := json.Marshal(blk)
	h := sha256.New()
	h.Write(j)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (bc blockchain) chainValid() bool {
	if len(bc.chain) < 2 {
		return true
	}
	for cindex := 1; cindex < len(bc.chain); cindex++ {
		prevBlock := bc.chain[cindex-1]
		curBlock := bc.chain[cindex]
		if curBlock.PreviousHash != prevBlock.Hash() {
			return false
		}
		prevProof := prevBlock.Proof
		curProof := curBlock.Proof
		h := sha256.New()
		h.Write([]byte(strconv.Itoa(curProof*curProof - prevProof*prevProof)))
		bs := h.Sum(nil)
		if !strings.HasPrefix(fmt.Sprintf("%x", bs), "0000") {
			return false
		}
	}
	return true
}

func init() {
	bc = blockchain{}
}

func main() {
	http.HandleFunc("/mine_block", mineBlockHandler)
	http.ListenAndServe(":8080", nil)
}

func mineBlockHandler(w http.ResponseWriter, r *http.Request) {

}

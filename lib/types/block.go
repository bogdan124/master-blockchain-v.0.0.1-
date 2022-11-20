package types

import (
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
)

//block struct
type Block struct {
	Index      int      `json:"index"`
	Timestamp  int64    `json:"timestamp"`
	Tx         [][]byte `json:"tx"`
	MerkleRoot []byte   `json:merkle`
	PrevHash   string   `json:"prevHash"`
	Hash       string   `json:"hash"`
	Nonce      int      `json:"nonce"`
	Reward     int      `json:"reward"`
	Coinbase   int      `json:"coinbase"`
	Miner      string   `json:"miner"`
}

//transaction struct
type Transaction struct {
	BlockNumber int
	Time        int64
	Hash        string
	Inputs      []TxInput
	Outputs     []TxOutput
}

//CalculateHash hashes the values of a TestContent
func (t Transaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.Hash + string(t.Time))); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t Transaction) Equals(other merkletree.Content) (bool, error) {
	return t.Hash == other.(Transaction).Hash, nil
}

type TxInput struct {
	Txid      int64
	Value     int64
	Signature string
	PubKey    string
}

type TxOutput struct {
	Txid       int64
	Value      int64
	PubKeyHash string
	Signature  string
}

type UTXO struct {
	Txid      string
	Index     int
	PubKey    string
	Signature string
	Value     int
}

type UTXOSet struct {
	UTXOs []UTXO
}

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

//wallet struct
type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

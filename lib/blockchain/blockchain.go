package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	//include from types block struct
	"github.com/cbergoon/merkletree"

	"blockchain-demo/lib/cryptogeneration"
	"blockchain-demo/lib/fileop"
	"blockchain-demo/lib/merkle"
	"blockchain-demo/lib/types"
	"blockchain-demo/lib/utils"
)

func CreateGenesisBlock(addressTo string) types.Block {
	//var reward int64
	//reward = 100
	//get private and public key

	publicKey, privateKey := cryptogeneration.GetPublicPrivateKeys()
	//create genesis block
	genesisBlock := types.Block{
		Index:      1,
		Timestamp:  time.Now().Unix(),
		Tx:         nil,
		PrevHash:   "",
		Hash:       "",
		MerkleRoot: nil,
		Nonce:      0,
		Reward:     100,
		Coinbase:   100,
		Miner:      string(publicKey),
	}
	//get private and public key
	//_, privateKey := cryptogeneration.GetPublicPrivateKeys()
	//add to transaction input

	bytesTx, _ := json.Marshal(CreateTransaction(privateKey, addressTo, nil, 100, 0, 1))
	genesisBlock.Tx = append(genesisBlock.Tx, bytesTx)
	//convert block to json
	//json, _ := json.Marshal(genesisBlock)
	//write to levelDB
	//convert int to byte
	//index := utils.IntToBytes(genesisBlock.Index)

	//fileop.PutInDB("db/blocks/blockchain", []byte(index), json)

	return genesisBlock
}

func CreateBlock() types.Block {
	//get last block
	lastkey := fileop.GetLastKey("db/blocks/blockchain")
	lastBlock := fileop.GetFromDB("db/blocks/blockchain", lastkey)
	blockDB := types.Block{}
	json.Unmarshal(lastBlock, &blockDB)

	//get public and private key
	publicKey, _ := cryptogeneration.GetPublicPrivateKeys()
	//create new block
	newBlock := types.Block{
		Index:      blockDB.Index + 1,
		Timestamp:  time.Now().Unix(),
		PrevHash:   blockDB.Hash,
		Hash:       "Genesis",
		MerkleRoot: nil,
		Nonce:      0,
		Reward:     100,
		Coinbase:   blockDB.Coinbase + 100,
		Miner:      string(publicKey),
	}

	return newBlock
}

func ProofOfWork(address string) types.Block {
	getLastBlockKey := fileop.GetLastKey("db/blocks/blockchain")
	if len(getLastBlockKey) == 0 {
		if getLastBlockKey == nil {
			//fmt.Println("Genesis:", string(getLastBlockKey))
			//create genesis block
			return CreateGenesisBlock(address)
		}
	}
	newBlock := CreateBlock()
	//get private and public key
	_, privateKey := cryptogeneration.GetPublicPrivateKeys()
	//create list of transaction
	var tx []merkletree.Content
	//get all transaction in levelDB
	transaction := fileop.GetAllKeys("db/mempool/valide")
	if len(transaction) != 0 {
		//unmarshall transaction
		for _, v := range transaction {
			data := fileop.GetFromDB("db/mempool/valide", v)
			transactionDB := types.Transaction{}
			json.Unmarshal(data, &transactionDB)
			transactionDB.BlockNumber = newBlock.Index
			//add to list of transaction
			tx = append(tx, transactionDB)
		}
	}

	tx = append(tx, CreateTransaction(privateKey, address, nil, 100, 0, 0))

	//create merkle tree
	merkleTree, _ := merkle.BuildTree(tx)
	listOfTransaction, RootNode := merkle.ExportLeafs(merkleTree)

	newBlock.Tx = listOfTransaction
	newBlock.MerkleRoot = RootNode
	//reconstruct merkle tree from string
	//merkleTree, _ := merkle.ReconstructTree(newBlock.Tx)
	//get merkle tree from string

	iterator := 0
	for {
		newBlock.Nonce = iterator
		storeHash := CalculateHash(newBlock)
		if storeHash[:2] == "00" {
			//fmt.Println("Hash:", storeHash)
			newBlock.Hash = storeHash
			//clear db//mempool/valide
			fileop.EraseAllKeys("db/mempool/valide")
			//fmt.Println("proof_of_work_func:", newBlock)
			return newBlock
			break
		}
		iterator += 1
	}

	return types.Block{}
}

//put valid block in levelDB
func PutBlockBlockchain(bl types.Block) {
	//convert block to json
	blockchainBytes, err := json.Marshal(bl)
	if err != nil {
		fmt.Println("error:", err)
	}
	//convert int to byte
	index := utils.IntToBytes(bl.Index)
	//write to levelDB
	fileop.PutInDB("db/blocks/blockchain", []byte("block"+string(index)), blockchainBytes)

}

//put valid block in levelDB
func PutBlockInvalid(bl types.Block) {
	//convert block to json
	json, _ := json.Marshal(bl)
	//write to levelDB
	fileop.PutInDB("db/blocks/invalid", []byte(bl.Hash), json)
}

//erase block invalid
func EraseBlockInvalidLs() {

	//write to levelDB
	fileop.EraseAllKeys("db/blocks/invalid")
}

func CalculateHash(bl types.Block) string {
	//   Calculate the hash of a block
	//create a string with the block data
	data := fmt.Sprintf("%d%d%s%d%d", bl.Index, bl.Timestamp, bl.Tx, bl.PrevHash, bl.Nonce)
	//calculate the hash of the block data
	hash := sha256.Sum256([]byte(data))
	//return the hash
	return hex.EncodeToString(hash[:])
}

func ValidateBlock(bl types.Block) bool {

	//   Validate the hash of a block
	//read data file from blocks
	data_file, err := os.Open("blocks/data")
	if err != nil {
		fmt.Println(err)
	}
	//read file
	fileInfo, err := data_file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	//create a slice of bytes with the size of the file
	data := make([]byte, fileInfo.Size())
	//read the file
	data_file.Read(data)
	//print the file
	//fmt.Println(string(data))

	//read last block file
	file, err := os.Open(string(data))
	if err != nil {

		fmt.Println(err)
	}

	//read file
	fileInfoBlock, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	//put block inside a struct of type block
	var block types.Block
	//create a slice of bytes with the size of the file
	dataBlock := make([]byte, fileInfoBlock.Size())
	//read the file
	file.Read(dataBlock)
	//marshal bytes to json
	json.Unmarshal(dataBlock, &block)

	defer file.Close()
	block.Index = bl.Index
	block.PrevHash = bl.PrevHash
	block.Timestamp = bl.Timestamp
	block.Reward = bl.Reward
	block.Coinbase = bl.Coinbase
	block.Nonce = bl.Nonce
	block.Tx = bl.Tx
	block.Miner = bl.Miner
	block.Hash = ""

	//calculate the hash of the block
	hash := CalculateHash(block)

	fmt.Println("====Block is validated====")
	fmt.Println("Hash: ++", bl, "++")
	fmt.Println("-", block, "-")

	//print the hash
	fmt.Println(hash)
	//print the hash of the block
	fmt.Println(bl.Hash)
	//compare the hash of the block with the hash calculated
	if hash == bl.Hash {
		return true
	}
	return false

}

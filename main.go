package main

//import golelevel
import (
	"blockchain-demo/lib/blockchain"
	"blockchain-demo/lib/cryptogeneration"
	"blockchain-demo/lib/fileop"
	"blockchain-demo/lib/network"
	"blockchain-demo/lib/types"
	"blockchain-demo/lib/utils"
	"encoding/json"
	"fmt"
)

func main() {
	fileop.EraseAllKeys("db/peers")

	///testFile()
	network.RunSourceNode()
	//testKeySign()
	//testBlockchain()
	//test()
	//testKeys()

}

func testKeySign() {
	//generate public private key
	pub, priv := cryptogeneration.CreatePairPublicPrivateKey()
	//sign a message
	signature := cryptogeneration.SingUsingKey(priv, pub, []byte("asdasdadasdasdsa"))
	//verify signature
	if cryptogeneration.VerifySign(signature, pub, []byte("asdasdadasdasdsa")) {
		fmt.Println("Signature is valid")
	} else {
		fmt.Println("Signature is not valid")
	}

}

func testFile() {

	//get address
	pub, _ := cryptogeneration.GetPublicPrivateKeys()
	//add 10
	for i := 0; i < 10; i++ {
		blk := blockchain.ProofOfWork(string(pub))
		blockchain.PutBlockBlockchain(blk)
	}
	//get all keys
	//
	keys := fileop.GetAllKeys("db/blocks/blockchain")
	fmt.Println("Keys: ", keys)
	//iterate over list and print it
	for _, key := range keys {
		data := fileop.GetFromDB("db/blocks/blockchain", key)
		var block types.Block
		json.Unmarshal(data, &block)
		fmt.Println("________________________________")
		fmt.Println("Key: ", string(key))
		fmt.Println("Blocknumber: ", block.Index)
		fmt.Println("________________________________")

	}
	//fileop.GetLastKeyBeta("db/blocks/blockchain")
	//for i := 0; i < 4; i++ {

	//clear screen
	//convert struct to bytes
	/*	blkockBytesSec, err := json.Marshal(blk)
		if err != nil {
			fmt.Println("error converting struct to bytes")
		}*/
	//add \n at blkockbytes
	//blkockBytes := []byte(string(blkockBytesSec) + "\n")
	//println("||===========================")
	//println(blkockBytes)
	//println(string(blkockBytes))
	//println("||===========================")

	println("||===========================")
	//}
	println("????????===========================")
	//get all keys
	keys = fileop.GetAllKeys("db/blocks/blockchain")
	fmt.Println("Keys: ", keys)
	//iterate over list and print it
	for _, key := range keys {
		data := fileop.GetFromDB("db/blocks/blockchain", key)
		var block types.Block
		json.Unmarshal(data, &block)
		fmt.Println("________________________________")
		fmt.Println("Blocknumber: ", block.Index)
		fmt.Println("________________________________")

	}

}

func testBlockchain() {
	pub, _ := cryptogeneration.CreatePairPublicPrivateKey()
	blockchain.CreateGenesisBlock(string(pub))
	//convert int to byte
	var number int

	b := utils.IntToBytes(number)

	//get a block index from db
	block := fileop.GetFromDB("db/blocks/valid", b)
	//unmarshal block
	var block1 types.Block
	json.Unmarshal(block, &block1)
	fmt.Println("Block: ", block1)
	fmt.Println("==========")
	blockchain.ProofOfWork(string(pub))
	fileop.EraseAllKeys("db/blocks/valid")

}

func testKeys() {
	//create a key pair
	pub, priv := cryptogeneration.CreatePairPublicPrivateKey()
	//fmt.Println("Public key: ", pub)
	//	fmt.Println("Private key: ", priv)
	//read from leveldb
	//wallet1
	var wallet1 []byte
	var wallet types.Wallet
	wallet1 = []byte("wallet1")

	//get wallet from levedb
	walletBytes := fileop.GetFromDB("db/usr", wallet1)

	//unmarshal wallet
	json.Unmarshal(walletBytes, &wallet)
	//fmt.Println("Wallet: ", string(wallet.PrivateKey))
	//fmt.Println("Wallet: ", string(wallet.PublicKey))
	//fmt.Println("Wallet: ", string(priv))
	//fmt.Println("Public key: ", privsec)
	//compare privsec and priv
	if string(wallet.PrivateKey) == string(priv) {
		fmt.Println("Keys are the same")
	} else {
		fmt.Println("Keys are not the same")
	}
	fmt.Println(cryptogeneration.GetPrivateKeyAndValidatePublicKey(pub))
	//d, b := cryptogeneration.GetPublicPrivateKeys()
	//fmt.Println("Public key: ", string(d))
	//fmt.Println("Private key: ", string(b))
	fileop.EraseAllKeys("db/usr")

}

func test() {
	fileop.CreateDB()
	//put into db peers
	d1 := []byte("data")
	fileop.PutInDB("db/peers", d1, d1)
	//fileop.EraseAllKeys("db/peers")
	fileop.GetLastKey("db/peers")
	//delete
	//DeleteFromDB("db/peers", d1)
	keys := fileop.GetAllKeys("db/peers")
	fmt.Println(string(keys[0]))
	//update values
	d2 := []byte("data32")
	fileop.UpdateValue("db/peers", d1, d2)
	//get all keys
	keys = fileop.GetAllKeys("db/peers")
	fmt.Println(string(keys[0]))

	//get from db peers
	d := fileop.GetFromDB("db/peers", []byte("data"))
	fmt.Println(d)
	//erase all keys
	fileop.EraseAllKeys("db/peers")
	d = fileop.GetFromDB("db/peers", []byte("data"))
	fmt.Println(d)

	fmt.Println(string(d))
}

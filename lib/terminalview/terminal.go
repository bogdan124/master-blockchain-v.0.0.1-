package terminalview

import (
	"blockchain-demo/lib/blockchain"
	"blockchain-demo/lib/cryptogeneration"
	"blockchain-demo/lib/fileop"
	"blockchain-demo/lib/merkle"
	"blockchain-demo/lib/types"
	"blockchain-demo/lib/utils"
	"encoding/json"
	"strconv"

	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p/core/host"
)

//function that takes input about diffrent options available
func TerminalView(node host.Host) {
	//print a menu with 8 options
	//1. Create a new wallet
	//2. View your wallet
	//3. Send coins
	//4. View your transactions
	//5. Mine
	//6. Options
	//7. Exit

	theEntireMenu()

	//read data from keyboard
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	//remove \n from text
	text = text[:len(text)-1]
	//verify the input
	switch text {
	case "1":
		createWallet()
	case "2":
		viewWallet()
	case "3":
		sendCoins(node)
	case "4":
		viewTransactions()
	case "5":
		mine(node)
	case "6":
		options()
	case "7":
		exit()
	case "8":
		addPeers(node)
	}

}

func theEntireMenu() {
	fmt.Println("+--------------------+")
	fmt.Println("|1.)Create new wallet|")
	fmt.Println("+--------------------+")
	fmt.Println("|2.)View your wallet |")
	fmt.Println("+--------------------+")
	fmt.Println("|3).Send coins       |")
	fmt.Println("+--------------------+")
	fmt.Println("|4.)View transactions|")
	fmt.Println("+--------------------+")
	fmt.Println("|5.)Mine             |")
	fmt.Println("+--------------------+")
	fmt.Println("|6.)Options          |")
	fmt.Println("+--------------------+")
	fmt.Println("|7.)Exit             |")
	fmt.Println("+--------------------+")
	fmt.Println("|8.)Add peers        |")
	fmt.Println("+--------------------+")
}

func createWallet() {
	fmt.Println("+-------------------+")
	fmt.Println("|Create a new wallet|")
	fmt.Println("+-------------------+")
	_, pubString := cryptogeneration.CreatePairPublicPrivateKey()
	fmt.Println("Your wallet has been created")
	fmt.Println("Public key: ", pubString)
	fmt.Println("check file generated")

}

func viewWallet() {
	fmt.Println("+----------------+")
	fmt.Println("|View your wallet|")
	fmt.Println("+----------------+")
	//get wallet from file
	publicKey, _ := cryptogeneration.GetPublicPrivateKeys()
	fmt.Println("Your wallet: ", string(publicKey))
	data := fileop.GetAllKeys("db/wallet/utxo")
	fmt.Println("Your balance: ", data)
	for _, utxo := range data {
		var utxoStruct types.UTXO
		json.Unmarshal(utxo, &utxoStruct)
		fmt.Println("Index: ", utxoStruct.Index)
		fmt.Println("PubKey: ", utxoStruct.PubKey)
		fmt.Println("Txid: ", utxoStruct.Txid)
		fmt.Println("Value: ", utxoStruct.Value)
	}

}

func sendCoins(node host.Host) {
	fmt.Println("+----------------+")
	fmt.Println("|Send coins      |")
	fmt.Println("+----------------+")

	//read data from keyboard
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("connect peer addres: ")
	peer, _ := reader.ReadString('\n')

	//remove \n from text
	peer = peer[:len(peer)-1]

	//add to money to transaction
	//read data from keyboard
	reader = bufio.NewReader(os.Stdin)
	fmt.Print("amount: ")
	amount, _ := reader.ReadString('\n')

	//add destination wallet to transaction
	//read data from keyboard
	reader = bufio.NewReader(os.Stdin)
	fmt.Print("destination wallet: ")
	walletDest, _ := reader.ReadString('\n')

	//read your public key from file
	//get public and private keys
	publicKey, privateKey := cryptogeneration.GetPublicPrivateKeys()

	//create transaction
	//convert string to int64
	amountInt, err := strconv.ParseInt(amount[:len(amount)-1], 10, 64)
	if err != nil {
		fmt.Println("error converting string to int64")
	}
	//include CreateTransaction from transaction
	transaction := blockchain.CreateTransaction(privateKey, walletDest, publicKey, amountInt, 0, 0)
	//convert struct to bytes
	tx, _ := json.Marshal(transaction)
	//add \n at blkockbytes
	tx = []byte(string(tx) + "\n")

	//send to all peers
	utils.SendToAllPeers(node, tx, 1)

}

func viewTransactions() {
	fmt.Println("+-----------------+")
	fmt.Println("|View transactions|")
	fmt.Println("+-----------------+")

	//read valide transaction db and display
	tx_valide_pool := fileop.GetAllKeys("db/mempool/valide")

	for _, tx := range tx_valide_pool {
		//convert bytes to struct
		var transaction types.Transaction
		json.Unmarshal(tx, &transaction)
		fmt.Println("Blocknumber: ", transaction.BlockNumber)
		fmt.Println("Hash: ", transaction.Hash)
		fmt.Println("Inputs: ", transaction.Inputs)
		fmt.Println("Outputs: ", transaction.Outputs)
	}
}

func mine(node host.Host) {
	fmt.Println("+----------------+")
	fmt.Println("|Mine            |")
	fmt.Println("+----------------+")

	//get address
	pub, _ := cryptogeneration.GetPublicPrivateKeys()

	blk := blockchain.ProofOfWork(string(pub))
	//clear screen
	//convert struct to bytes
	blkockBytesSec, err := json.Marshal(blk)
	if err != nil {
		fmt.Println("error converting struct to bytes")
	}
	//add \n at blkockbytes
	blkockBytes := []byte(string(blkockBytesSec) + "\n")

	blockchain.PutBlockBlockchain(blk)
	//fileop.PutInDB("db/blocks/blockchain", []byte(blk.Hash), blkockBytesSec)
	bt := fileop.GetLastKey("db/blocks/blockchain")
	data := fileop.GetFromDB("db/blocks/blockchain", bt)
	var block types.Block
	json.Unmarshal(data, &block)
	fmt.Println("__________________________________")
	fmt.Println("|Blocknumber: ", block.Index)
	fmt.Println("|Hash: ", block.Hash)
	fmt.Println("|Miner: ", block.Miner)
	fmt.Println("|Timestamp: ", block.Timestamp)
	fmt.Println("|Nonce: ", block.Nonce)

	//get public key
	publicKey, _ := cryptogeneration.GetPublicPrivateKeys()

	fmt.Println("+++++++++++++++++++++++++++++++++++++")
	fmt.Println("Transactions in block: ")
	fmt.Println("+++++++++++++++++++++++++++++++++++++")
	merkle.RebuildTree(blk.Tx, blk.MerkleRoot)
	merkle.PrintLeafs(blk.Tx)

	//get all inputs tx
	for _, tx := range blk.Tx {
		//convert bytes to struct
		var transactionData types.Transaction
		json.Unmarshal(tx, &transactionData)
		//get all TX
		blockchain.AddYourTransactionToWallet(transactionData, string(publicKey))
	}

	utils.SendToAllPeers(node, blkockBytes, 0)

}

func options() {
	fmt.Println("+----------------+")
	fmt.Println("|Options         |")
	fmt.Println("+----------------+")

}

func exit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func addPeers(node host.Host) {
	fmt.Println("+----------------+")
	fmt.Println("|Add peers       |")
	fmt.Println("+----------------+")

	reader := bufio.NewReader(os.Stdin)
	//read int from keyboards
	fmt.Print("number of peers: ")
	peers, _ := reader.ReadString('\n')
	//remove \n from text
	peers = peers[:len(peers)-1]
	//convert string to int
	peersInt, _ := strconv.Atoi(peers)

	//add number of peersInt
	for i := 0; i < peersInt; i++ {
		//read data from keyboard
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("peer addres: ")
		peer, _ := reader.ReadString('\n')
		var available string
		available = "available"
		//check if peer is alread in db
		utils.AddPeersToDB(peer, []byte(available))
		//remove \n from text
		peer = peer[:len(peer)-1]

	}

}

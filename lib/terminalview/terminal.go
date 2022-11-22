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

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/ttacon/chalk"
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
	//clear screen clscr from os
	//clear screen
	fmt.Print("\033[H\033[2J")
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

	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+------------------------------+"))
	fmt.Println(lime("|      Welcome to the menu     |"))
	fmt.Println(lime("+----------------------------- +"))
	fmt.Println(lime("|1. Create a new wallet        |"))
	fmt.Println(lime("|2. View your wallet           |"))
	fmt.Println(lime("|3. Send coins                 |"))
	fmt.Println(lime("|4. View your transactions     |"))
	fmt.Println(lime("|5. Mine                       |"))
	fmt.Println(lime("|6. Options                    |"))
	fmt.Println(lime("|7. Exit                       |"))
	fmt.Println(lime("|8. Add peers                  |"))
	fmt.Println(lime("+------------------------------+"))
}

func createWallet() {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+-------------------+"))
	fmt.Println(lime("|Create a new wallet|"))
	fmt.Println(lime("+-------------------+"))
	_, pubString := cryptogeneration.CreatePairPublicPrivateKey()
	fmt.Println(lime("Your wallet has been created"))
	fmt.Println(lime("Public key: "), pubString)
	fmt.Println(lime("check file generated"))

}

func viewWallet() {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+----------------+"))
	fmt.Println(lime("|View your wallet|"))
	fmt.Println(lime("+----------------+"))
	//get wallet from file
	publicKey, _ := cryptogeneration.GetPublicPrivateKeys()
	fmt.Println(lime("Your wallet: "), string(publicKey))
	data := fileop.GetAllKeys("db/wallet/utxo")
	var balance int64

	//get balance from utxo
	for _, value := range data {
		utxoBytes := fileop.GetFromDB("db/wallet/utxo", value)
		//convert bytes to struct
		var transaction types.Transaction
		json.Unmarshal(utxoBytes, &transaction)
		//get the amount
		for _, input := range transaction.Inputs {
			//convert to float
			balance += input.Value
		}
	}
	//convert balance to string
	balanceString := strconv.FormatInt(balance, 10)
	fmt.Println(lime("+-------------------------------------------+"))
	fmt.Println(lime("|Your balance: "), balanceString)
	fmt.Println(lime("--------------------------------------------+"))
	//print transaction
	for _, value := range data {
		utxoBytes := fileop.GetFromDB("db/wallet/utxo", value)
		fmt.Println(lime("Tx: "), string(utxoBytes))
		fmt.Println(lime("----------------------------------------------------------"))
	}

}

func sendCoins(node host.Host) {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	//red warning
	red := chalk.Red.NewStyle().
		WithTextStyle(chalk.Bold).
		Style

	fmt.Println(lime("+----------------+"))
	fmt.Println(lime("|Send coins      |"))
	fmt.Println(lime("+----------------+"))

	//add to money to transaction
	//read data from keyboard
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(lime("|amount: "))
	amount, _ := reader.ReadString('\n')

	//add destination wallet to transaction
	//read data from keyboard
	reader = bufio.NewReader(os.Stdin)
	fmt.Print(lime("|destination wallet: "))
	walletDest, _ := reader.ReadString('\n')

	//read your public key from file
	//get public and private keys
	publicKey, privateKey := cryptogeneration.GetPublicPrivateKeys()

	//create transaction
	//convert string to int64
	amountInt, err := strconv.ParseInt(amount[:len(amount)-1], 10, 64)
	if err != nil {
		fmt.Println(red("|error converting string to int64"))
	}
	//include CreateTransaction from transaction
	transaction := blockchain.CreateTransaction(privateKey, walletDest, publicKey, amountInt, 0, 0)
	//convert struct to bytes
	tx, _ := json.Marshal(transaction)
	//add \n at blkockbytes
	tx = []byte(string(tx) + "\n")
	fmt.Println(lime("|Transaction created: "), string(tx))
	//print line
	fmt.Println(lime("----------------------------------------------------------"))
	//send to all peers
	utils.SendToAllPeers(node, tx, 1)

}

func viewTransactions() {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+-----------------+"))
	fmt.Println(lime("|View transactions|"))
	fmt.Println(lime("+-----------------+"))

	//read valide transaction db and display
	tx_valide_pool := fileop.GetAllKeys("db/mempool/valide")

	for _, tx := range tx_valide_pool {
		//convert bytes to struct
		var transaction types.Transaction
		json.Unmarshal(tx, &transaction)
		fmt.Println(lime("|Blocknumber: "), transaction.BlockNumber)
		fmt.Println(lime("|Hash: "), transaction.Hash)
		fmt.Println(lime("|Inputs: "), transaction.Inputs)
		fmt.Println(lime("|Outputs: "), transaction.Outputs)
	}
}

func mine(node host.Host) {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
		//red warning
	red := chalk.Red.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+----------------+"))
	fmt.Println(lime("|Mine            |"))
	fmt.Println(lime("+----------------+"))

	//get address
	pub, _ := cryptogeneration.GetPublicPrivateKeys()

	blk := blockchain.ProofOfWork(string(pub))
	//clear screen
	//convert struct to bytes
	blkockBytesSec, err := json.Marshal(blk)
	if err != nil {
		fmt.Println(red("|error converting struct to bytes"))
	}
	//add \n at blkockbytes
	blkockBytes := []byte(string(blkockBytesSec) + "\n")

	blockchain.PutBlockBlockchain(blk)
	//fileop.PutInDB("db/blocks/blockchain", []byte(blk.Hash), blkockBytesSec)
	bt := fileop.GetLastKey("db/blocks/blockchain")
	data := fileop.GetFromDB("db/blocks/blockchain", bt)
	var block types.Block
	json.Unmarshal(data, &block)
	fmt.Println(lime("+-----------------------------------------------------------"))
	fmt.Println(lime("|Blocknumber: "), block.Index)
	fmt.Println(lime("|Hash: "), block.Hash)
	fmt.Println(lime("|Miner: "), block.Miner)
	fmt.Println(lime("|Timestamp: "), block.Timestamp)
	fmt.Println(lime("|Nonce: "), block.Nonce)

	//get public key
	publicKey, _ := cryptogeneration.GetPublicPrivateKeys()

	fmt.Println(lime("+++++++++++++++++++++++++++++++++++++"))
	fmt.Println(lime("Transactions in block: "))
	fmt.Println(lime("+++++++++++++++++++++++++++++++++++++"))
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
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+----------------+"))
	fmt.Println(lime("|Options         |"))
	fmt.Println(lime("+----------------+"))

}

func exit() {
	//red warning
	red := chalk.Red.NewStyle().
		WithTextStyle(chalk.Bold).
		Style

	fmt.Println(red("+---------------------------------+"))
	fmt.Println(red("|Exiting...                       |"))
	fmt.Println(red("|Received signal, shutting down...|"))
	fmt.Println(red("+---------------------------------+"))
	os.Exit(0)
}

func addPeers(node host.Host) {
	lime := chalk.Green.NewStyle().
		WithTextStyle(chalk.Bold).
		Style
	fmt.Println(lime("+----------------+"))
	fmt.Println(lime("|Add peers       |"))
	fmt.Println(lime("+----------------+"))

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

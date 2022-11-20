package blockchain

import (
	"blockchain-demo/lib/cryptogeneration"
	"blockchain-demo/lib/fileop"
	"blockchain-demo/lib/types"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

func CreateTransaction(privateKey []byte, addressTo string, addressFrom []byte, value int64, indexTxInput int64, indexTxOutput int64) types.Transaction {
	//create a new transaction
	//convert int 64 to byte
	bnumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(bnumber, uint64(value))
	//validate trasnsaction
	var transaction types.Transaction
	//transaction from a block
	if addressFrom == nil {
		//saveSig := cryptogeneration.SingUsingKey(privateKey, addressFrom, bnumber)
		outputTx := []types.TxOutput{}
		inputTx := []types.TxInput{}

		transaction = types.Transaction{
			BlockNumber: 0,
			Time:        time.Now().Unix(),
			Hash:        "",
			Inputs:      inputTx,
			Outputs:     outputTx,
		}
		transaction.Inputs = append(transaction.Inputs, types.TxInput{
			Txid:      0,
			Value:     value,
			Signature: "",
			PubKey:    "coinbase",
		})
		saveSig := cryptogeneration.SingUsingKey(privateKey, []byte(addressTo), bnumber)
		transaction.Outputs = append(transaction.Outputs, types.TxOutput{
			Txid:       0,
			Value:      value,
			PubKeyHash: addressTo,
			Signature:  string(saveSig),
		})
		jsonTransaction, _ := json.Marshal(transaction)
		transaction.Hash = cryptogeneration.GenerateSha256(jsonTransaction)
		return transaction
	} else {
		//transaction from a wallet
		saveSig := cryptogeneration.SingUsingKey(privateKey, []byte(addressFrom), bnumber)
		outputTx := []types.TxOutput{}
		inputTx := []types.TxInput{}

		transaction = types.Transaction{
			BlockNumber: 0,
			Time:        time.Now().Unix(),
			Hash:        "",
			Inputs:      inputTx,
			Outputs:     outputTx,
		}

		transaction.Inputs = GetPastTransactions(addressFrom)
		balance := GetBalanceFromInputsForWallet(addressFrom)
		if balance < value {
			fmt.Println("Not enough balance")
			return types.Transaction{}
		} else {

			//value to send
			transaction.Outputs = append(transaction.Outputs, types.TxOutput{
				Txid:       indexTxOutput,
				Value:      value,
				PubKeyHash: string(addressFrom),
				Signature:  string(saveSig),
			})
			//value to return to your wallet
			transaction.Outputs = append(transaction.Outputs, types.TxOutput{
				Txid:       indexTxOutput,
				Value:      balance - value,
				PubKeyHash: addressTo,
				Signature:  string(saveSig),
			})
		}
		//json transaction
		jsonTransaction, _ := json.Marshal(transaction)
		transaction.Hash = cryptogeneration.GenerateSha256(jsonTransaction)
		return transaction
	}
	//put transaction into badgerdb
	//PutTransaction(transaction)

}

func ValidateTransaction(tx types.Transaction) bool {
	//check if the transaction output singature is valid
	for _, output := range tx.Outputs {
		//validate signature
		if !cryptogeneration.VerifySign([]byte(output.PubKeyHash), []byte(output.Signature), []byte(output.PubKeyHash)) {
			fmt.Println("Invalid signature")
			return false
		}
	}

	return true
}

func AddMemPoolTransaction(transaction types.Transaction) {

	//add to bytes Transaction
	jsonTransaction, _ := json.Marshal(transaction)
	fileop.PutInDB("db/mempool/valide", []byte(transaction.Hash), jsonTransaction)
}

func AddYourTransactionToWallet(tx types.Transaction, walletPubKey string) {
	//check if the transaction input your walletpubkey is equel with output pubkey
	fmt.Println("wallet", tx.Outputs)
	for _, out := range tx.Outputs {
		fmt.Println("Public key :", out.PubKeyHash)
		if out.PubKeyHash == walletPubKey {
			//marshal transaction to json
			jsonTransaction, _ := json.Marshal(tx)
			fmt.Println("jsonTransaction", jsonTransaction)
			//add transaction to wallet usr
			fileop.PutInDB("db/wallet/utxo", []byte(tx.Hash), jsonTransaction)
		}
	}
}

func GetPastTransactions(address []byte) []types.TxInput {
	//get all transaction from badgerdb
	//return a list of transaction
	//read db/wallet/utxo from db
	keys := fileop.GetAllKeys("db/wallet/utxo")
	if len(keys) == 0 {
		return nil
	}
	//list of transaction
	var tx []types.TxInput
	//read all transaction
	for _, key := range keys {
		//read transaction
		transaction := fileop.GetFromDB("db/wallet/utxo", key)
		//unmarshal transaction
		var txInput types.TxInput
		json.Unmarshal(transaction, &txInput)
		//append transaction to list of transaction
		tx = append(tx, txInput)
	}
	//return list of transaction
	return tx
}

func GetBalanceFromInputsForWallet(address []byte) int64 {
	//get all transaction from badgerdb
	//return the balance
	tx := GetPastTransactions(address)
	if tx == nil {
		return 0
	}
	//balance
	var balance int64
	//read all transaction
	for _, input := range tx {
		//read transaction
		//unmarshal transaction
		//add value to balance
		balance += input.Value
	}
	//return balance
	return balance
}

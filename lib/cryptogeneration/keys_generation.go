package cryptogeneration

import (
	"blockchain-demo/lib/fileop"
	"blockchain-demo/lib/types"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	//crypto from standard library

	cryptolib "crypto"

	"github.com/libp2p/go-libp2p-core/crypto"
)

//function to create a public and private key
func CreateKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
	//   Create a private and public key pair
	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, nil, err
	}

	return priv, pub, nil
}

func CreatePairPublicPrivateKey() ([]byte, []byte) {

	//create a private and public key
	priv, pub, err := CreateKeyPair()
	if err != nil {
		panic(err)
	}

	//convert private key to bytes
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		panic(err)
	}

	//convert public key to bytes
	pubBytes, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		panic(err)
	}

	//convert private key to string
	//privString := string(privBytes)
	//convert public key to string and format it
	//pubString := FormatPublicAddress(pubBytes)
	//list of bytes to string
	privString := hex.EncodeToString(privBytes)
	pubString := hex.EncodeToString(pubBytes)

	//create a wallet variable
	wallet := types.Wallet{
		PrivateKey: []byte(privString),
		PublicKey:  []byte(pubString),
	}
	//convert wallet to bytes marshal
	walletBytes, _ := json.Marshal(wallet)
	//fmt.Println("walletBytes", walletBytes)
	//put wallet in the db
	fileop.PutInDB("db/usr", []byte("wallet1"), walletBytes)

	return []byte(pubString), []byte(privString)

}

func FormatPublicAddress(pubBytes []byte) string {
	//generate sha256 for public key
	pubHash := GenerateSha256(pubBytes)
	//print public key hash
	log.Printf("Public key hash: %s", pubHash)
	//concatenate 0x in front characters and remove 10 from hash\
	pubHashRet := "0xadd" + pubHash

	return pubHashRet
}

func ReformatPublicAddress(publickey string) string {
	//remove 0xadd from public key
	publickey = publickey[2:]
	return publickey
}

func GetPrivateKeyAndValidatePublicKey(publicKey []byte) bool {
	//get private key from db
	//aloc a variable
	var firstWallet []byte
	//get the wallet from db
	firstWallet = []byte("wallet1")

	//get wallet from levedb
	walletBytes := fileop.GetFromDB("db/usr", firstWallet)
	//fmt.Println("walletBytes", walletBytes)
	if walletBytes == nil {
		return false
	}
	//unmarshal wallet
	var wallet1 types.Wallet
	json.Unmarshal(walletBytes, &wallet1)
	//encode to bytes
	privKeyBytes, _ := hex.DecodeString(string(wallet1.PrivateKey))
	//unmarshal private key
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		log.Printf("Error unmarshalling private key: %s", err)
		return false
	}

	//get public key from private key
	pubKey := privKey.GetPublic()

	//convert public key to bytes
	pubBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		log.Printf("Error marshalling public key: %s", err)
		return false
	}
	//convert public key to string
	pubString := hex.EncodeToString(pubBytes)
	//convert public key to string
	//pubString2 := hex.EncodeToString(publicKey)

	//compare public keys
	//fmt.Println("pubString", pubString)
	//fmt.Println("publicKey", string(pubString2))
	if string(publicKey) != pubString {
		log.Printf("Public keys are not the same")
		return false
	}

	return true
}

func GenerateSha256(value []byte) string {
	//apllly standard library to create a hash sha256
	hash := cryptolib.SHA256.New()
	//write privKey to hash
	hash.Write(value)
	//get the hash
	hashed := hash.Sum(nil)
	//convert hash to string
	hashedString := hex.EncodeToString(hashed)
	return hashedString
}

func SingUsingKey(privateKey []byte, publicKey []byte, data []byte) []byte {
	if GetPrivateKeyAndValidatePublicKey(publicKey) {
		//decode private key
		privKeyBytes, _ := hex.DecodeString(string(privateKey))
		//unmarshal private key
		privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
		if err != nil {
			log.Printf("Error unmarshalling private key: %s", err)
			return nil
		}

		//sign data
		signature, err := privKey.Sign(data)
		if err != nil {
			log.Printf("Error signing data: %s", err)
			return nil
		}

		return signature

	} else {
		fmt.Println("Public key does not match private key")
		return nil
	}

}

func VerifySign(signature []byte, publicKey []byte, data []byte) bool {

	pubKey, err := crypto.UnmarshalPublicKey(publicKey)
	if err != nil {
		log.Printf("Error unmarshalling public key: %s", err)
		return false
	}

	//verify signature
	verify, err := pubKey.Verify(data, signature)
	if err != nil {
		return false
	}

	return verify
}

func GetPublicPrivateKeys() ([]byte, []byte) {
	//get public and private keys from db
	//aloc a variable
	var firstWallet []byte
	//get the wallet from db
	firstWallet = []byte("wallet1")

	//get wallet from levedb
	walletBytes := fileop.GetFromDB("db/usr", firstWallet)
	//fmt.Println("walletBytes", walletBytes)
	if walletBytes == nil {
		return nil, nil
	}
	//unmarshal wallet
	var wallet1 types.Wallet
	json.Unmarshal(walletBytes, &wallet1)

	return wallet1.PublicKey, wallet1.PrivateKey
}

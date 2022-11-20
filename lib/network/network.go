package network

import (
	"blockchain-demo/lib/blockchain"
	"blockchain-demo/lib/cryptogeneration"
	"blockchain-demo/lib/terminalview"
	"blockchain-demo/lib/utils"

	"blockchain-demo/lib/types"

	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func ReadTransactionProtocol(s network.Stream) error {
	//   Read the stream and print its content
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	connection := s.Conn()
	message = message[:len(message)-1]
	messageBytes := []byte(message)
	var transaction types.Transaction
	err = json.Unmarshal(messageBytes, &transaction)
	if err != nil {
		return nil
	}
	//if transaction is not null
	if transaction.Outputs != nil {

		//check if transaction is valid
		if blockchain.ValidateTransaction(transaction) {
			//get your public key
			pub, _ := cryptogeneration.GetPublicPrivateKeys()

			blockchain.AddYourTransactionToWallet(transaction, string(pub))
			blockchain.AddMemPoolTransaction(transaction)

		}
	} else {
		fmt.Println("Transaction is null")
	}

	log.Printf("Message from '%s': %s", connection.RemotePeer().String(), message)
	return nil
}

func ReadMessagesProtocol(s network.Stream) error {
	//   Read the stream and print its content
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	connection := s.Conn()

	log.Printf("Message from '%s': %s", connection.RemotePeer().String(), message)
	return nil
}

func ReadMineProtocol(s network.Stream) error {
	buf := bufio.NewReader(s)
	//read bytes from stream until EOF
	//buf.Size()
	message, err := buf.ReadString('\n')
	if err != nil {
		return err
	}
	//convert remove \n chars
	message = message[:len(message)-1]
	//convert string to bytes
	messageBytes := []byte(message)
	//convert to Block
	var block types.Block
	err = json.Unmarshal(messageBytes, &block)
	//print block
	if err != nil {
		return err
	}

	connection := s.Conn()
	//fmt.Println(block)
	//for the simple I will not validate the block I will just Add to Blockchain
	blockchain.PutBlockBlockchain(block)
	//print block
	fmt.Println("=========================================")
	fmt.Println("|      Block added to Blockchain        |")
	fmt.Println("=========================================")
	fmt.Println("|Index: ", block.Index)
	fmt.Println("|Coinbase: ", block.Coinbase)
	//convert merkleroot to hex
	fmt.Println("|MerkleRoot: ", string(block.MerkleRoot))
	fmt.Println("|Miner: ", block.Miner)
	fmt.Println("|Transactions: ", block.Tx)
	fmt.Println("|PrevHash: ", block.PrevHash)
	fmt.Println("|Hash: ", block.Hash)
	fmt.Println("|Nonce: ", block.Nonce)
	fmt.Println("|Timestamp: ", block.Timestamp)
	fmt.Println("=========================================")

	var createIncomingAdd string

	createIncomingAdd = connection.RemoteMultiaddr().String() + "/p2p/" + connection.RemotePeer().String() + "\n"

	//add to list of peers
	utils.AddPeersToDB(createIncomingAdd, []byte("available"))

	return nil
}

func RunSourceNode() {
	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
	node := CreateNode("/ip4/127.0.0.1/tcp/0")

	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])
	go LogicNodeInteraction(node)

	for {
		terminalview.TerminalView(node)
	}

}

func LogicNodeInteraction(node host.Host) {
	//Set stream handler for the "/hello/1.0.0" protocol
	go node.SetStreamHandler("/transaction/1.0.0", func(s network.Stream) {
		log.Printf("/transaction/1.0.0 stream created")
		err := ReadTransactionProtocol(s)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	go node.SetStreamHandler("/mine/1.0.0", func(s network.Stream) {
		log.Printf("/mine/1.0.0 stream created")
		err := ReadMineProtocol(s)
		if err != nil {
			log.Printf("Error: %s", err)
			s.Reset()
		} else {
			log.Printf("Closing stream")
			s.Close()
		}
	})

	go node.SetStreamHandler("/messages/1.0.0", func(s network.Stream) {
		log.Printf("/messages/1.0.0 stream created")
		err := ReadMessagesProtocol(s)
		if err != nil {
			log.Printf("Error: %s", err)
			s.Reset()
		} else {
			log.Printf("Closing stream")
			s.Close()
		}
	})

}

func CreateNode(hostip string) host.Host {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings(hostip),
		libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	//d := cryptogeneration.GetPrivateKeyFromFileAndValidatePublicKey("0xadd613596b13c8695004309e960f31a3e596d4bfb3d457e5f716076c8c8c5df8d0e")
	//print d
	//log.Printf("Private key: %s", d)
	//createPairPublicPrivateKey
	//createPairPublicPrivateKey("hello")
	return node
}

package utils

import (
	"blockchain-demo/lib/fileop"
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func OpenConnectionTransaction(peer_addres string, node host.Host, data []byte) {
	addr, err := multiaddr.NewMultiaddr(peer_addres)
	if err != nil {
		println("Error: address format wrong")
		return
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		println("Error: address format wrong")
		return
	}

	if err := node.Connect(context.Background(), *peer); err != nil {
		println("Error: connection failed")
		return
	}
	transaction, err := node.NewStream(context.Background(), peer.ID, "/transaction/1.0.0")
	if err != nil {
		//print err
		log.Fatal(err)
		println("Error: stream creation failed")
		return
	}

	log.Printf("Sending message...")

	//send string to stream
	_, err = transaction.Write(data)
	if err != nil {
		println("Error: message sending failed")
		return
	}

}

func OpenConnectionMine(peer_addres string, node host.Host, data []byte) {

	addr, err := multiaddr.NewMultiaddr(peer_addres)
	if err != nil {
		println("Error: address format wrong1")
		return
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		println("Error: address format wrong2")
		return
	}

	if err := node.Connect(context.Background(), *peer); err != nil {
		println("Error: connection failed")
		return
	}
	//we send  a block in this section that was mined
	stream, err := node.NewStream(context.Background(), peer.ID, "/mine/1.0.0")
	if err != nil {
		//print err
		log.Fatal(err)
		println("Error: stream creation failed")
		return
	}
	fmt.Print("Sending message...")

	//send string to stream
	_, err = stream.Write(data)

	if err != nil {
		println("Error: message sending failed")
		return
	}
}

func IntToBytes(number int) []byte {
	index := make([]byte, 8)
	binary.LittleEndian.PutUint64(index, uint64(number))
	return index
}

func SendToAllPeers(node host.Host, data []byte, type_to_send int) {
	//read from db all keys
	number := fileop.GetNumberOfKeys("db/peers")
	if number == 0 {
		fmt.Println()
	} else {

		//get all keys from db/peers
		keys := fileop.GetAllKeys("db/peers")
		fmt.Println("=================")
		//for each key send message
		for i := 0; i < len(keys); i++ {
			fmt.Println(string(keys[i]))
		}
		fmt.Println("================")
		//for each key convert to string
		for i := 0; i < number; i++ {
			peer_addres := string(keys[i])
			//remove \n from string
			peer_addres = peer_addres[:len(peer_addres)-1]
			fmt.Println("Sending to peer: ", peer_addres)
			if type_to_send == 0 {
				//send to blocks channel
				OpenConnectionMine(peer_addres, node, data)
			} else if type_to_send == 1 {
				//send to transaction channel
				OpenConnectionTransaction(peer_addres, node, data)
			}
		}
	}
}

func AddPeersToDB(peers string, available []byte) {
	//add peer to db
	fileop.PutInDB("db/peers", []byte(peers), available)
}

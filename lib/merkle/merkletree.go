package merkle

import (
	"blockchain-demo/lib/types"
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cbergoon/merkletree"
)

/*
//Equals tests for equality of two Contents
func (t types.Transaction) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}
*/

func BuildTree(list []merkletree.Content) (*merkletree.MerkleTree, string) {

	//Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()

	mrHex := fmt.Sprintf("%x", mr)

	return t, mrHex
}

func ExportLeafs(m *merkletree.MerkleTree) ([][]byte, []byte) {

	var expLeafs [][]byte
	//var g *merkletree.Node
	for _, l := range m.Leafs {
		//convert l.C to bytes
		b, err := json.Marshal(l.C)
		if err != nil {
			log.Fatal(err)
		}
		//pus to expLeafs
		expLeafs = append(expLeafs, b)
	}
	//remove dublicates
	expLeafs = removeDuplicates(expLeafs)

	return expLeafs, m.MerkleRoot()
}

func removeDuplicates(elements [][]byte) [][]byte {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := [][]byte{}

	for v := range elements {
		if encountered[string(elements[v])] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[string(elements[v])] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func RebuildTree(leafs [][]byte, roottrre []byte) bool {
	//rebuild tree
	var list []merkletree.Content
	for _, v := range leafs {
		var t types.Transaction
		err := json.Unmarshal(v, &t)
		if err != nil {
			log.Fatal(err)
			return false
		}
		list = append(list, t)
	}
	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
		return false
	}
	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	//check roottrre
	if bytes.Equal(mr, roottrre) {
		log.Println("Roots are equal")
		return true
	} else {
		log.Println("Roots are not equal")
		return false
	}

}

//print transanction from leaf
func PrintLeafs(leafs [][]byte) {
	for _, v := range leafs {
		var t types.Transaction
		err := json.Unmarshal(v, &t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BlockNumber:", t.BlockNumber)
		fmt.Println("Hash:", t.Hash)
		fmt.Println("Time:", t.Time)
		fmt.Println("===Inputs===")
		for _, v := range t.Inputs {
			fmt.Println("Txid:", v.Txid)
			fmt.Println("Value:", v.Value)
			fmt.Println("Signature:", v.Signature, len(v.Signature), "|")
			fmt.Println("PubKey:", v.PubKey)
		}
		fmt.Println("===Outputs===")
		for _, v := range t.Outputs {
			fmt.Println("Txid:", v.Txid)
			fmt.Println("Value:", v.Value)
			fmt.Println("PubKeyHash:", v.PubKeyHash)
			fmt.Println("Signature:", v.Signature, "|")
		}
		fmt.Println("________________________________________")
	}
}

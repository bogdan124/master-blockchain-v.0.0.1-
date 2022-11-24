package fileop

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func CreateDB() bool {
	//create a new database
	//here we store blocks that are valid and ready to put into blockchain
	d1, err := leveldb.OpenFile("db/blocks/valid", nil)
	if err != nil {
		return false
	}
	defer d1.Close()
	//here we store blocks that come from the peer and are not validated
	d2, err5 := leveldb.OpenFile("db/blocks/invalid", nil)
	if err5 != nil {
		return false
	}
	defer d2.Close()
	//here we store blocks that are  valid and are ready to put into blockchain
	d3, err6 := leveldb.OpenFile("db/blocks/blockchain", nil)
	if err6 != nil {
		return false
	}
	defer d3.Close()
	//here are store all the transaction the will be put in a block
	d4, err3 := leveldb.OpenFile("db/mempool/valide", nil)
	if err3 != nil {
		return false
	}
	defer d4.Close()
	//here we store all the transaction that are not valid because
	//they came from the network and need to wait for validation
	d5, err4 := leveldb.OpenFile("db/mempool/invalide", nil)
	if err4 != nil {
		return false
	}
	defer d5.Close()
	//here we store all the past peers because we need to send message towards them later
	d6, err2 := leveldb.OpenFile("db/peers", nil)
	if err2 != nil {
		return false
	}
	defer d6.Close()

	db7, err7 := leveldb.OpenFile("db/usr", nil)
	if err7 != nil {
		return false
	}
	defer db7.Close()

	db8, err8 := leveldb.OpenFile("db/wallet/utxo", nil)
	if err8 != nil {
		return false
	}
	defer db8.Close()

	return true
}

//give a leght of bytes and put it into a db
func PutInDB(dbname string, key []byte, value []byte) bool {
	//open database
	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		//print err
		fmt.Println(err)
		return false
	}

	err = db.Put(key, value, nil)

	//close database
	defer db.Close()
	return err != nil
	//return true
}

//get a value from a db
func GetFromDB(dbname string, key []byte) []byte {
	o := &opt.Options{
		ReadOnly: true,
	}
	//open database
	db, err := leveldb.OpenFile(dbname, o)
	if err != nil {
		return nil
	}

	data, err := db.Get(key, nil)
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	//close database
	defer db.Close()
	return data
}

func DeleteFromDB(dbname string, key []byte) bool {
	//open database
	db, _ := leveldb.OpenFile(dbname, nil)

	err := db.Delete(key, nil)
	//close database
	defer db.Close()
	return err != nil
}

//get all the keys from a db
func GetAllKeys(dbname string) [][]byte {
	o := &opt.Options{
		ReadOnly: true,
	}
	//open database
	db, err := leveldb.OpenFile(dbname, o)
	if err != nil {
		return nil
	}

	iter := db.NewIterator(nil, nil)
	var keys [][]byte
	for iter.Next() {
		key := iter.Key()
		//copy key
		keycopy := make([]byte, len(key))
		copy(keycopy, key)
		keys = append(keys, keycopy)
	}
	iter.Release()

	//close database
	defer db.Close()

	return keys
}

//erase all the keys from a db
func EraseAllKeys(dbname string) bool {
	//open database
	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		return false
	}

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		err := db.Delete(iter.Key(), nil)
		if err != nil {
			return false
		}
	}
	iter.Release()
	//close database
	defer db.Close()
	return err != nil
}

//update a value read
func UpdateValue(dbname string, key []byte, value []byte) bool {
	//open database
	db, _ := leveldb.OpenFile(dbname, nil)

	err := db.Put(key, value, nil)
	//close database
	defer db.Close()
	return err != nil
}

//get last key from the db
func GetLastKey(dbname string) []byte {
	o := &opt.Options{
		ReadOnly: true,
	}
	//open database
	db, err := leveldb.OpenFile(dbname, o)
	if err != nil {
		return nil
	}

	iter := db.NewIterator(nil, nil)
	var lastkey []byte
	for iter.Next() {
		lastkey = iter.Key()
	}
	iter.Release()

	//close database
	defer db.Close()
	return lastkey
}
func GetLastKeyBeta(dbname string) []byte {
	o := &opt.Options{
		ReadOnly: true,
	}
	//open database
	db, err := leveldb.OpenFile(dbname, o)
	if err != nil {
		return nil
	}
	iter := db.NewIterator(nil, nil)

	var key []byte
	for ok := iter.Seek(key); ok; ok = iter.Next() {
		key = iter.Key()
		//val := iter.Value()
		//fmt.Printf("key=%s, value=%s", string(key), string(val))
	}

	iter.Release() // Note: you should first get data and then release iterator
	err = iter.Error()
	return key
}

func GetNumberOfKeys(dbname string) int {
	o := &opt.Options{
		ReadOnly: true,
	}
	//open database
	db, _ := leveldb.OpenFile(dbname, o)

	iter := db.NewIterator(nil, nil)
	number := 0
	for iter.Next() {
		number++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return 0
	}
	//close database
	defer db.Close()
	return number
}

func PutMultipleValuesDB(dbname string, list_keys [][]byte, list_data [][]byte) {
	//open database
	db, _ := leveldb.OpenFile(dbname, nil)

	batch := new(leveldb.Batch)
	for i := 0; i < len(list_data); i++ {
		batch.Put(list_keys[i], list_data[i])
	}
	db.Write(batch, nil)

	//close database
	defer db.Close()
}

/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"time"
//	"encoding/json"
	"strings"
	"strconv"
	"io"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

//NRxPLLBMIJ
//jx52X9ASHP
//posARbsxeV

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	Use(args)
	//if len(args) != 2 {
	//	return shim.Error("Incorrect arguments. Expecting a key and a value")
	//}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	stub.PutState("state", []byte("0"))
	stub.PutState("user1", []byte(""))
	stub.PutState("user2", []byte(""))
	stub.PutState("user3", []byte(""))
	stub.PutState("dealer", []byte(""))
	stub.PutState("Deadline", []byte{})
//	if err != nil {
//		return shim.Error(fmt.Sprintf("Failed to create asset"))
//	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	/*if fn == "set" {
		result, err = set(stub, args)
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}*/
	if fn == "set_deadline" {
		result, err = set_deadline(stub, args) //set the dealine
	} else if fn == "set_state" {
		result, err = set_state(stub, args) //when the time is up, change the state
	} else if fn == "view_state" {
		result, err = view_state(stub, args) //show the current state
//	} else if fn == "view_sign" {
//		result, err = view_sign(stub, args)
	} else if fn == "send_deal" {
		result, err = send_deal(stub, args) //show the current state
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

func send_deal(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	result, err := set_state(stub, args)
	flag := 0
	if strings.Compare(string(result), "Deadline due") == 0 {
		return "Deadline due, subscription failed", nil
	}
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a value")
	}
	user := 0
	key := ""
	if strings.Compare(args[0], "user1") == 0 {
		user = 1
		key = "NRxPLLBMIJ"
	}
	if strings.Compare(args[0], "user2") == 0 {
		user = 2
		key = "jx52X9ASHP"
	}
	if strings.Compare(args[0], "user3") == 0 {
		user = 3
		key = "posARbsxeV"
	}
	if user == 0 {
		return "User not found", nil
	}
	if flag == 1 {
		return "Deadline is due, update failed", nil
	}
	cur_dealstr := decrypt([]byte(key)[:16], args[1])
	cur_deal, err := strconv.Atoi(cur_dealstr)
	if err != nil {
		return "", fmt.Errorf("dealstr problem")
	}
	str, _ := stub.GetState(args[0])
	if string(str) != "" {
		return "", fmt.Errorf("Submission exists, aborting")
	}
	stub.PutState(args[0], []byte(cur_dealstr))
	str, _ = stub.GetState("dealer")
	if strings.Compare(string(str), "") == 0 {
		stub.PutState("dealer", []byte(args[0]))
			return "new dealer", nil
	} else {
		bestbyte, _ := stub.GetState(string(str))
		beststr := string(bestbyte)
		best, _ := strconv.Atoi(beststr)
		if best > cur_deal {
			stub.PutState("dealer", []byte(args[0]))
			return "dealer changed", nil
		} else {
			return "dealer remains", nil
		}
	}
}

func view_state(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	result, _ := set_state(stub, args)
	flag := 0
	if strings.Compare(string(result), "Deadline due") == 0 {
		flag = 1
	}
	if len(args) == 0 {
		str, err := stub.GetState("state")
		return string(str), err
	}
	if len(args) != 1 {//dealer / deal
		return "", fmt.Errorf("Incorrect arguments. Expecting a value")
	}
	if flag == 0 {
		return "Deadline is not due, information not unfolded", nil
	}
	if strings.Compare(string(args[0]), "dealer") == 0 {
		str, err := stub.GetState("dealer")
		return string(str), err
	} else if strings.Compare(string(args[0]), "deal") == 0 {
		str, err := stub.GetState("dealer")
		str, err = stub.GetState(string(str))
		return string(str), err
	} else {
		return "View state: no such value", nil
	}
}

func set_deadline(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {//year month date hour minute
		return "", fmt.Errorf("Incorrect arguments. Expecting Year, Month, Day, Hour, Minute")
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	month ,err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	day ,err := strconv.Atoi(args[2])
	if err != nil {
		return "", err
	}
	hour ,err := strconv.Atoi(args[3])
	if err != nil {
		return "", err
	}
	minute ,err := strconv.Atoi(args[4])
	if err != nil {
		return "", err
	}
	if year > 2020 || year < 2000 {
		return "", fmt.Errorf("Year error")
	}
	if month > 12 || month < 1 {
		return "", fmt.Errorf("Month error")
	}
	if day > 31 || day < 1 {
		return "", fmt.Errorf("Day error")
	}
	if hour > 23 || hour < 0 {
		return "", fmt.Errorf("Hour error")
	}
	if minute > 59 || minute < 0 {
		return "", fmt.Errorf("Min error")
	}
	dl := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
	bytedl, err := dl.MarshalJSON()
	if err != nil {
		return "", err
	}
	err = stub.PutState("Deadline", bytedl)
	if err != nil {
		return "", err
	}
	return "Set deadline", err
}

func set_state(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	bytedl, err := stub.GetState("Deadline")
	if err != nil {
		return "fail to get time", err
	}
	var dl time.Time
	err = dl.UnmarshalJSON(bytedl)
	if err != nil {
		return "", err
	}
	now := time.Now()
	if now.After(dl) {
		err := stub.PutState("state", []byte("1"))
		if err != nil {
			return "fail to update state", err
		}
		return "Deadline due", err
	} else {
		return "Deadline not due", err
	}
}

l
// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func Use(vals ...interface{}){
	for _, val := range vals {
		_ = val
	}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"strings"
	"strconv"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	err = stub.PutState("SNUM", []byte(strconv.Itoa(0)))
	err = stub.PutState("UNUM", []byte(strconv.Itoa(0)))
	err = stub.PutState("CANDIDATE_1", []byte(strconv.Itoa(0)))
	err = stub.PutState("CANDIDATE_2", []byte(strconv.Itoa(0)))
	err = stub.PutState("CANDIDATE_3", []byte(strconv.Itoa(0)))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
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
	if fn == "verify" {
		result, err = verify(stub, args)
	} else if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "pray" {
		result, err = pray(stub, args)
	} else if fn == "minmin" {
		result, err = minmin(stub, args)
	} else if fn == "addadd" {
		result, err = addadd(stub, args)
	} else if fn == "adduser" {
		result, err = adduser(stub, args)
	} else if fn == "votefor" {
		result, err = votefor(stub, args)
	} else if fn == "seepoll" {
		result, err = seepoll(stub, args)
	} else {// assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

func seepoll(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 { //uname, passwd
		return "", fmt.Errorf("SEEPOLL FORMAT ERROR")
	}
	flag := 0
	if strings.Compare(string(args[0]), "CANDIDATE_1") == 0 {
		flag = 1
	}
	if strings.Compare(string(args[0]), "CANDIDATE_2") == 0 {
		flag = 1
	}
	if strings.Compare(string(args[0]), "CANDIDATE_3") == 0 {
		flag = 1
	}
	if flag == 0 {
		return "", fmt.Errorf("CANDIDATE NF")
	}
	value, err := stub.GetState(args[0])
	return string(value), err
}

func verify(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 { //uname, passwd
		return "", fmt.Errorf("VERIFY FORMAT ERROR")
	}
	s, err := stub.GetState("UNUM")
	unum, err := strconv.Atoi(string(s))
	if err != nil {
		return "", fmt.Errorf("failed to get UNUM with error: %s", err)
	}
	i := 0
	for i < unum {
		s, _ = stub.GetState("UID" + strconv.Itoa(i))
		p, _ := stub.GetState("PID" + strconv.Itoa(i))
		t, _ := stub.GetState("SID" + strconv.Itoa(i))
		if strings.Compare(string(s), args[0]) == 0 {
			if strings.Compare(string(p), args[1]) != 0 {
				return "PASSWD WRONG", err
			}
			if strings.Compare(string(t), "true") == 0 {
				return "VOTE USED" , err
			} else {
				return "VOTE NOT USED", err
			}
		}
		i++
	}
	return "NO SUCH USER", err
}

func votefor(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 { //uname, passwd, candid
		return "", fmt.Errorf("NOTEFOR FORMAT ERROR")
	}
	s, err := stub.GetState("UNUM")
	unum, err := strconv.Atoi(string(s))
	if err != nil {
		return "", fmt.Errorf("failed to get UNUM with error: %s", err)
	}
	i := 0
	for i < unum {
		s, _ = stub.GetState("UID" + strconv.Itoa(i))
		p, _ := stub.GetState("PID" + strconv.Itoa(i))
		t, _ := stub.GetState("SID" + strconv.Itoa(i))
		if strings.Compare(string(s), args[0]) == 0 {
			if strings.Compare(string(p), args[1]) != 0 {
				return "", fmt.Errorf("PASSWD WRONG")
			}
			if strings.Compare(string(t), "false") != 0 {
				return "", fmt.Errorf("ACCOUNT USED")
			}
			flag := 0
			if strings.Compare(string(args[2]), "CANDIDATE_1") == 0 {
				flag = 1
			}
			if strings.Compare(string(args[2]), "CANDIDATE_2") == 0 {
				flag = 1
			}
			if strings.Compare(string(args[2]), "CANDIDATE_3") == 0 {
				flag = 1
			}
			if flag == 0 {
				return "", fmt.Errorf("CANDIDATE NF")
			}
			s, err := stub.GetState(args[2])
			value, err := strconv.Atoi(string(s))
			str := strconv.Itoa(value + 1)
			err = stub.PutState(args[2], []byte(str))
			err = stub.PutState("SID" + strconv.Itoa(i), []byte("true"))
			return "VOTED FOR " + args[2], err
		}
		i++
	}
	return "", fmt.Errorf("USER NF")
}

func adduser(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("ADDUSER FORMAT ERROR")
	}
	s, err := stub.GetState("UNUM")
	unum, err := strconv.Atoi(string(s))
	if err != nil {
		return "", fmt.Errorf("failed to get UNUM with error: %s", err)
	}
	i := 0
	for i < unum {
		s, err = stub.GetState("UID" + strconv.Itoa(i))
		if strings.Compare(string(s), args[0]) == 0 {
			return "", fmt.Errorf("USER EXISTS")
		}
		i++
	}
	err = stub.PutState("UID" + strconv.Itoa(unum), []byte(args[0]))
	err = stub.PutState("PID" + strconv.Itoa(unum), []byte(args[1]))
	err = stub.PutState("SID" + strconv.Itoa(unum), []byte("false"))
	unum++
	str := strconv.Itoa(unum)
	err = stub.PutState("UNUM", []byte(str))
	return "ADD USER " + args[0], nil
}

//##################################################

func pray(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("No, just pray")
	}
	fmt.Printf("Success, Jesus received your prayer")
	return "true", nil
}

func addadd(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("No args")
	}
	s, err := stub.GetState("SNUM")
	value, err := strconv.Atoi(string(s))
	if err != nil {
		return "", fmt.Errorf("failed to get snum with error: %s", err)
	}
//	if value == nil {
//		return "", fmt.Errorf("SNUM NF")
//	}
	if value >= 10 {
		return "", fmt.Errorf("SNUM is at ceiling")
	}
	str := strconv.Itoa(value + 1)
	err = stub.PutState("SNUM", []byte(str))
	if err != nil {
		return "", fmt.Errorf("Failed to perform addadd")
	}
	return str, nil
}

func minmin(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("No args")
	}
	s, err := stub.GetState("SNUM")
	value, err := strconv.Atoi(string(s))
	if err != nil {
		return "", fmt.Errorf("Failed to get SNUM with error: %s", err)
	}
//	if value == nil {
//		return "", fmt.Errorf("SNUM NF")
//	}
	if value <= 0 {
		return "", fmt.Errorf("SNUM is at bottom")
	}
	str := strconv.Itoa(value - 1)
	err = stub.PutState("SNUM", []byte(str))
	if err != nil {
		return "", fmt.Errorf("Failed to perform minmin")
	}
	return str, nil
}
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

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

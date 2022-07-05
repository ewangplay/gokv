package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// KVStore provides functions for managing Key/Value
type KVStore struct {
	contractapi.Contract
}

// InitLedger ...
func (s *KVStore) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// Set adds a new Key/Value pair to the world state
func (s *KVStore) Set(ctx contractapi.TransactionContextInterface, k string, v string) error {
	if len(k) == 0 {
		return fmt.Errorf("input param key is missing")
	}
	if len(v) == 0 {
		return fmt.Errorf("input param value is missing")
	}
	return ctx.GetStub().PutState(k, []byte(v))
}

// Get returns the value stored in the world state with given key
func (s *KVStore) Get(ctx contractapi.TransactionContextInterface, k string) (string, error) {
	if len(k) == 0 {
		return "", fmt.Errorf("input param key is missing")
	}
	v, err := ctx.GetStub().GetState(k)
	if err != nil {
		return "", err
	}
	if len(v) == 0 {
		return "", fmt.Errorf("%s does not exist", k)
	}
	return string(v), nil
}

// Delete deletes the Key / Value from the world state with given key
func (s *KVStore) Delete(ctx contractapi.TransactionContextInterface, k string) error {
	if len(k) == 0 {
		return fmt.Errorf("input param key is missing")
	}
	return ctx.GetStub().DelState(k)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(KVStore))
	if err != nil {
		fmt.Printf("Error create chaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v\n", err)
	}
}

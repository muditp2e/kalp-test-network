/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chaincode-go/chaincode"
	"log"

	"github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

/*
contract := kalpsdk.Contract{IsPayableContract: false}
    contract.Logger = kalpsdk.NewLogger()
    nftChaincode, err := kalpsdk.NewChaincode(&kalpCBDC.SmartContract{Contract: contract})
*/

func main() {
	contract := kalpsdk.Contract{IsPayableContract: false}
	contract.Logger = kalpsdk.NewLogger()
	nftChaincode, err := kalpsdk.NewChaincode(&chaincode.SmartContract{Contract: contract})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := nftChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}

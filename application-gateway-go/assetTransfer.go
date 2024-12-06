/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID        = "p2epro"
	cryptoPath   = "/home/mudit/Documents/projects/kamalp2e/kalp-test-network/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	channelName := "kalptantra"
	chaincodeName := "gini1"

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.

	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	commands := []string{"initLedger", "initialize", "mint", "balanceOf", "totalSupply", "name", "symbol", "decimals", "allow", "deny", "approve", "allowance", "cleanAll"}
	var s string
	for i, fn := range commands {
		s = s + fmt.Sprintf("%d: %s\n", i, fn)
	}
	for {
		// Take user input
		fmt.Println("Enter a command:\n" + s)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		// input := "3"

		// Trim whitespace and convert to lowercase
		commandIndex, err := strconv.Atoi(strings.TrimSpace(input))

		if err != nil || commandIndex >= len(commands) {
			fmt.Println("Invalid command index:", input)
			return
		}

		switch commandIndex {
		case 0:
			initLedger(contract)
		case 1:
			initialize(contract, "GINI", "GINI", "klp-6b616c70627269646775-cc")
		case 2:
			mint(contract, "user1", "100000000")
		case 3:
			balanceOf(contract, "user1")
			// balanceOf(contract, "klp-6b616c70627269646775-cc")
			balanceOf(contract, "klp-6b616c70627269646766-cc")
		case 4:
			GenericEvaluateTransaction(contract, "TotalSupply")
		case 5:
			GenericEvaluateTransaction(contract, "Name")
		case 6:
			GenericEvaluateTransaction(contract, "Symbol")
		case 7:
			GenericEvaluateTransaction(contract, "Decimals")
		case 8:
			GenericSubmitTransaction(contract, "Allow", "878ab9087t756c7857d596876e78687087f898")
		case 9:
			GenericSubmitTransaction(contract, "Deny", "878ab9087t756c7857d596876e78687087f898")
		case 10:
			GenericSubmitTransaction(contract, "Approve", "klp-6b616c70627269646766-cc", "1234567")
		case 11:
			GenericEvaluateTransaction(contract, "Allowance", "user1", "klp-6b616c70627269646766-cc")
		case 12:
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"docType":"UTXO"}}`)
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"docType":"UserRoleMap"}}`)
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"_id":{"$regex":"denyList"}}}`)
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"docType":"Approval"}}`)
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"_id":"name"}}`)
			GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"_id":"symbol"}}`)
			// GenericEvaluateTransaction(contract, "DeleteDocTypes", `{"selector":{"_id":"gasFees"}}`)

		}
	}

	// initLedger(contract)
	// initialize(contract, "GINI", "GINI", "klp-6b616c70627269646766-cc")
	//getAllAssets(contract)
	// createAsset(contract)
	// readAssetByID(contract)
	//transferAssetAsync(contract)
	//exampleErrorHandling(contract)
}

func initialize(contract *client.Contract, name, symbol, vestingContract string) {
	fmt.Printf("\n--> Submit Transaction: Initialize, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("Initialize", name, symbol, vestingContract)
	if err != nil {
		fmt.Println("failed to submit transaction: %w", getErrorDetails(err))
		return
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func mint(contract *client.Contract, address, amount string) {
	fmt.Printf("\n--> Submit Transaction: mint\n")

	_, err := contract.SubmitTransaction("mint", address, amount)
	if err != nil {
		fmt.Println("failed to submit transaction: %w", getErrorDetails(err))
		return
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func balanceOf(contract *client.Contract, owner string) {
	fmt.Println("\n--> Evaluate Transaction: balanceOf")

	evaluateResult, err := contract.EvaluateTransaction("BalanceOf", owner)
	if err != nil {
		fmt.Println("failed to evaluate transaction: %w", err)
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
func GenericEvaluateTransaction(contract *client.Contract, fnName string, args ...string) {
	fmt.Println("\n--> Evaluate Transaction: " + fnName)

	evaluateResult, err := contract.EvaluateTransaction(fnName, args...)
	if err != nil {
		fmt.Println("failed to evaluate transaction: %w", err)
	}
	// result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", evaluateResult)
}

func GenericSubmitTransaction(contract *client.Contract, fnName string, args ...string) {
	fmt.Println("\n--> Submit Transaction:" + fnName)

	_, err := contract.SubmitTransaction(fnName, args...)
	if err != nil {
		fmt.Println("failed to submit transaction: %w", getErrorDetails(err))
		return
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func extractErrorMessageFromGrpcStatus(transactionError *client.TransactionError) error {
	// log := utils.GetLoggingInstance().Log

	if transactionError != nil {
		grpcStatus := transactionError.GRPCStatus()
		if grpcStatus != nil {
			details := grpcStatus.Details()
			if len(details) > 0 {
				detail, ok := details[0].(*gateway.ErrorDetail)
				if ok {
					fmt.Println("Error details:", detail.Message)
					return fmt.Errorf("%v", detail.Message)
				}
			}
		}
	}
	return nil
}

func getErrorDetails(err error) error {
	// log := utils.GetLoggingInstance().Log

	transactionError, ok := err.(*client.TransactionError)
	if ok {
		err := extractErrorMessageFromGrpcStatus(transactionError)
		if err != nil {
			return err
		}
	}

	endorseError, ok := err.(*client.EndorseError)
	if ok {
		transactionError := endorseError.TransactionError
		err := extractErrorMessageFromGrpcStatus(transactionError)
		if err != nil {
			return err
		}
	}

	submitError, ok := err.(*client.SubmitError)
	if ok {
		transactionError := submitError.TransactionError
		err := extractErrorMessageFromGrpcStatus(transactionError)
		if err != nil {
			return err
		}
	}

	commitStatusError, ok := err.(*client.CommitStatusError)
	if ok {
		transactionError := commitStatusError.TransactionError
		err := extractErrorMessageFromGrpcStatus(transactionError)
		if err != nil {
			return err
		}
	}

	errorStatus := status.Convert(err)
	if errorStatus != nil {
		details := errorStatus.Details()
		if len(details) > 0 {
			detail := details[0]
			errorDetail, ok := detail.(*gateway.ErrorDetail)
			if ok {
				fmt.Println("Error details:", errorDetail.Message)
				return fmt.Errorf("%v", errorDetail.Message)
			}
		}
	}

	return err
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction to query ledger state.
func getAllAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner and AppraisedValue arguments \n")

	_, err := contract.SubmitTransaction("CreateKyc", "e99958785f5af2d9d637822a7ff09d537882c93c", "e99958785f5af2d9d637822a7ff09d537882c93c", "")
	// _, err := contract.SubmitTransaction("CreateKycV1", `{"docType":"kyc","userId":"8d4f5238b2bcf6ebdffe02db4c0ff876b5c815aa","kycId":"8d4f5238b2bcf6ebdffe02db4c0ff876b5c815aa","kycHash":"","isAbove18":true,"isAbove21":true,"isAbove60":false,"gender":"Male","kycProvider":"OneTrust","region":"APAC","country":"India"}`)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction by assetID to query ledger state.
func readAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	// evaluateResult, err := contract.EvaluateTransaction("KycExists", "08b2e895ea2149fb975fff3daff714e5bc27ae1e")
	evaluateResult, err := contract.EvaluateTransaction("KycExists", "e99958785f5af2d9d637822a7ff09d537882c93c")
	// evaluateResult, err := contract.EvaluateTransaction("IsUserAbove60", "8d4f5238b2bcf6ebdffe02db4c0ff876b5c815aa")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferAssetAsync(contract *client.Contract) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, updates existing asset owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(assetId, "Mark"))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
func exampleErrorHandling(contract *client.Contract) {
	fmt.Println("\n--> Submit Transaction: UpdateAsset asset70, asset70 does not exist and should return an error")

	_, err := contract.SubmitTransaction("UpdateAsset", "asset70", "blue", "5", "Tomoko", "300")
	if err == nil {
		panic("******** FAILED to return an error")
	}

	fmt.Println("*** Successfully caught the error:")

	switch err := err.(type) {
	case *client.EndorseError:
		fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.SubmitError:
		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
		}
	case *client.CommitError:
		fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
	default:
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
}

// Format JSON data
func formatJSON(data []byte) string {
	if len(data) == 0 {
		return "No data available"
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

package executor

import (
	"context"
	"fmt"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/tryoutbounder/golang-soroban/pkg/rpc/protocol"
)

func (e *Executor) SimulateContractCall(
	contractAddress xdr.ScAddress,
	sourceAccount txnbuild.Account,
	args []xdr.ScVal,
	functionName xdr.ScSymbol,
) (*xdr.ScVal, error) {

	transactionXdr, err := buildContractTx(contractAddress, sourceAccount, args, functionName)
	if err != nil {
		return nil, err
	}

	transactionBase64, err := transactionXdr.Base64()
	if err != nil {
		return nil, err
	}

	response, err := e.rpc.SimulateTransaction(
		context.TODO(),
		protocol.SimulateTransactionRequest{
			Transaction: transactionBase64,
		},
	)

	if err != nil {
		return nil, err
	}

	var responseScVal xdr.ScVal

	if len(response.Results) != 1 {
		return nil, fmt.Errorf("unexpected number of simulation results: %d", len(response.Results))
	}

	err = xdr.SafeUnmarshalBase64(
		*response.Results[0].ReturnValueXDR,
		&responseScVal,
	)

	if err != nil {
		return nil, err
	}

	return &responseScVal, err
}

// Untested
func (e *Executor) SubmitContractCall(
	contractAddress xdr.ScAddress,
	sourceAccount txnbuild.Account,
	args []xdr.ScVal,
	functionName xdr.ScSymbol,
	networkPassphrase string,
	signingKeypairs []*keypair.Full,
) (string, error) {
	fmt.Printf("Building contract transaction for contract address: %v\n", contractAddress)
	transactionXdr, err := buildContractTx(contractAddress, sourceAccount, args, functionName)
	if err != nil {
		fmt.Printf("Error building contract transaction: %v\n", err)
		return "", err
	}

	fmt.Printf("Signing transaction with %d keypairs\n", len(signingKeypairs))
	for _, keypair := range signingKeypairs {
		fmt.Printf("Signing with keypair: %s\n", keypair.Address())
		transactionXdr, err = transactionXdr.Sign(
			networkPassphrase,
			keypair,
		)
		if err != nil {
			fmt.Printf("Error signing transaction: %v\n", err)
			return "", err
		}
	}

	fmt.Println("Converting transaction to Base64")
	transactionBase64, err := transactionXdr.Base64()
	if err != nil {
		fmt.Printf("Error converting transaction to Base64: %v\n", err)
		return "", err
	}

	fmt.Println("Sending transaction")
	response, err := e.rpc.SendTransaction(
		context.TODO(),
		protocol.SendTransactionRequest{
			Transaction: transactionBase64,
		},
	)
	if err != nil {
		fmt.Printf("Error sending transaction: %v\n", err)
		return "", err
	}

	if response.ErrorResultXDR != "" {
		fmt.Printf("Received error result XDR: %s\n", response.ErrorResultXDR)
		var xdrErr xdr.ScError

		err := xdr.SafeUnmarshalBase64(response.ErrorResultXDR, &xdrErr)
		if err != nil {
			fmt.Printf("Error unmarshaling error result: %v\n", err)
			return "", err
		}
		fmt.Printf("Contract error: %+v\n", xdrErr)
		return "", fmt.Errorf("contract error: %+v", xdrErr)
	}
	fmt.Printf("Transaction sent successfully. Hash: %s\n", response.Hash)
	return response.Hash, nil
}

package helpers

import (
	"fmt"
	"math/big"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func ContractAddressToScAddress(tokenContractStr string) (xdr.ScAddress, error) {
	var contractAddress xdr.ScAddress
	tokenAddress, err := strkey.Decode(strkey.VersionByteContract, tokenContractStr)
	if err != nil {
		return contractAddress, fmt.Errorf("error decoding token contract: %v", err)
	}

	var tokenAddressHash xdr.ContractId
	copy(tokenAddressHash[:], tokenAddress)

	contractAddress, err = xdr.NewScAddress(xdr.ScAddressTypeScAddressTypeContract, tokenAddressHash)
	if err != nil {
		return contractAddress, fmt.Errorf("error creating contract address: %v", err)
	}

	return contractAddress, nil
}

func StellarAddressToScAddress(address string) (xdr.ScAddress, error) {
	var scAddress xdr.ScAddress
	// Decode the address string to get the raw bytes
	addressBytes, err := strkey.Decode(strkey.VersionByteAccountID, address)
	if err != nil {
		return scAddress, err
	}

	// Convert to Uint256 for the AccountId
	var addressUint256 xdr.Uint256
	copy(addressUint256[:], addressBytes)

	balanceAccount, err := xdr.NewAccountId(xdr.PublicKeyTypePublicKeyTypeEd25519, addressUint256)
	if err != nil {
		return scAddress, err
	}

	scAddress, err = xdr.NewScAddress(xdr.ScAddressTypeScAddressTypeAccount, balanceAccount)
	if err != nil {
		return scAddress, err
	}

	return scAddress, nil
}

func EncodeContractAddress(contractId xdr.ContractId) (string, error) {
	encodedStr, err := strkey.Encode(strkey.VersionByteContract, contractId[:])
	if err != nil {
		return "", fmt.Errorf("failed to encode address: %v", err)
	}
	return encodedStr, nil
}

// I128ToFloat64 converts an xdr.Int128Parts to a float64 with offset division
func I128ToFloat64(i128 xdr.Int128Parts, offset float64) float64 {
	// Create a big.Int from the high and low parts
	high := big.NewInt(int64(i128.Hi))
	low := big.NewInt(int64(i128.Lo))

	// Shift high part by 64 bits and add low part
	high.Lsh(high, 64)
	result := big.NewInt(0)
	result.Add(high, low)

	// Convert to float64
	f, _ := result.Float64()
	return f / offset
}

// I128ToInt64 converts an xdr.Int128Parts to an int64
func I128ToInt64(i128 xdr.Int128Parts) int64 {
	// Create a big.Int from the high and low parts
	high := big.NewInt(int64(i128.Hi))
	low := big.NewInt(int64(i128.Lo))

	// Shift high part by 64 bits and add low part
	high.Lsh(high, 64)
	result := big.NewInt(0)
	result.Add(high, low)

	// Convert to int64
	return result.Int64()
}

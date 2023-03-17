package main

import (
	"encoding/hex"
	"fmt" // TODO:  REMOVE BEFORE SUBMITTING
)

type encryptionOracle func(m []byte) (iv [2]byte, c []byte)

// Helper function to handle checking errors returned by functions. If err is
// not nil, Go aborts by calling panic()
func checkError(err error, msg string) {
	if err != nil {
		panic(err)
	}
}

// Helper to build arrays of bytes from hex strings
// (like those that you type when interacting with the
// router binary)
func hexStringToBytes(s string) []byte {
	bytes, err := hex.DecodeString(s)
	checkError(err, "Error decoding hex string to bytes")

	return bytes
}

func performAttack(oracle encryptionOracle, firstIv [2]byte, firstCiphertext []byte) [8]byte {
	// TODO:  perform your attack, calling oracle() to get a new IV, ciphertext pair
	// for a given message
	//
	// See main.go for an example of how this function is called.
	//
	// TODO:  Perform your attack!
	// Here are some examples of the inputs and how to use the oracle() functionconst
	//
	// **NOTE**:  Be sure to comment out any print statements before submitting,
	// as the "fmt" package is disallowed by the autograder for security reasons!!!

	fmt.Printf("First IV:  %x, First c:  %x\n", firstIv, firstCiphertext) // TODO:  REMOVE BEFORE SUBMITTING

	// Given m, get the next IV, ciphertext pair from the oracle (ie, the simulated router)
	m := hexStringToBytes("0000000000000000")
	iv, c := oracle(m)
	fmt.Printf("Got IV:  %x, First c:  %x\n", iv, c) // TODO:  REMOVE BEFORE SUBMITTING

	// When done, return the key as an array of 8 bytes
	key := make([]byte, 8)

	// If your key variable is a slice, here's an example of how to convert it to an array
	var keyAsArray [8]byte
	copy(keyAsArray[:], key)
	return keyAsArray
}

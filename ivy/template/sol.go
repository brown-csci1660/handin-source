package main

type encryptionOracle func(m []byte) (iv [2]byte, c []byte)

func performAttack(oracle encryptionOracle, firstIV [2]byte, firstCiphertext []byte) (key [8]byte) {
	// TODO: Implement your attack here!
	return
}

package main

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <key>\n", os.Args[0])
		os.Exit(1)
	}
	bytes, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse key as hex: %v\n", err)
		os.Exit(1)
	}
	if len(bytes) != 8 {
		fmt.Fprintln(os.Stderr, "key must be 8 bytes (16 hexadecimal digits)")
		os.Exit(1)
	}

	var k [8]byte
	copy(k[:], bytes)
	oracle := makeOracle(k)
	firstIV, firstCiphertext := oracle(k[:])

	fmt.Print("Running attack... ")
	t0 := time.Now()
	recoveredK := performAttack(oracle, firstIV, firstCiphertext)
	t1 := time.Now()
	fmt.Printf("completed in %v\n", t1.Sub(t0))
	if recoveredK == k {
		fmt.Println("Key successfully recovered.")
	} else {
		fmt.Println("ERROR: Wrong key recovered:")
		fmt.Printf("got  %v\n", hex.EncodeToString(recoveredK[:]))
		fmt.Printf("want %v\n", hex.EncodeToString(k[:]))
	}
}

// makeOracle returns an encryption oracle
// which will encrypt packets with the given key.
func makeOracle(k [8]byte) encryptionOracle {
	return func(m []byte) (iv [2]byte, c []byte) {
		// initialize iv
		iv[0] = byte(rand.Uint32())
		iv[1] = byte(rand.Uint32())

		// iv || k
		seed := make([]byte, 10)
		copy(seed, iv[:])
		copy(seed[2:], k[:])

		cipher, _ := rc4.NewCipher(seed)
		c = make([]byte, len(m))
		cipher.XORKeyStream(c, m)
		return iv, c
	}
}

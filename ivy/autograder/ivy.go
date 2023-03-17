package main

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const tries = 10

func main() {
	successes := 0
	for i := 0; i < tries; i++ {
		var k [8]byte
		binary.BigEndian.PutUint32(k[:], rand.Uint32())
		binary.BigEndian.PutUint32(k[4:], rand.Uint32())
		oracle := makeOracle(k)
		firstIV, firstCiphertext := oracle(k[:])

		recoveredK := performAttack(oracle, firstIV, firstCiphertext)
		if recoveredK == k {
			successes++
		}
	}

	ratio := float64(successes) / float64(tries)
	points := int(100.0 * ratio)
	comment := fmt.Sprintf("%v tries/%v successes", tries, successes)
	fmt.Printf("problem:ivy points:%v comment:%v\n", points, comment)
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

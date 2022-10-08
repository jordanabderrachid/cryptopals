package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println("hello")
}

func hexToBase64(in string) string {
	bytes, err := hex.DecodeString(in)
	panicIfErr(err)

	return base64.StdEncoding.EncodeToString(bytes)
}

func xor(lhs, rhs string) string {
	lhsBytes, err := hex.DecodeString(lhs)
	panicIfErr(err)

	rhsBytes, err := hex.DecodeString(rhs)
	panicIfErr(err)

	return hex.EncodeToString(xorBytes(lhsBytes, rhsBytes))
}

func xorBytes(lhs, rhs []byte) []byte {
	if len(lhs) != len(rhs) {
		panic("tried to xor []byte of unequal length")
	}

	size := len(rhs)
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		res[i] = lhs[i] ^ rhs[i]
	}
	return res
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

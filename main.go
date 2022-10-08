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

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

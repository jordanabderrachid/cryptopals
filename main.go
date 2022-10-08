package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
)

func main() {
	in := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	type candidate struct {
		key   byte
		plain string
		score float64
	}

	candidates := make([]candidate, 0)
	for i := 0; i < int(math.Pow(2, 8)); i++ {
		key := byte(i)
		encrypted, err := hex.DecodeString(in)
		panicIfErr(err)

		plainBuilder := &strings.Builder{}

		for _, b := range encrypted {
			plainBuilder.WriteByte(b ^ key)
		}

		plain := plainBuilder.String()
		candidates = append(candidates, candidate{
			key:   key,
			plain: plain,
			score: frequencyScore(plain),
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score < candidates[j].score
	})

	for _, c := range candidates {
		fmt.Printf("%s %f\n", c.plain, c.score)
	}
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

func frequencyScore(input string) float64 {
	// https://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
	type frequencyEntry struct {
		target string
		value  float64
	}
	entries := []frequencyEntry{
		{target: "e", value: 12.02},
		{target: "t", value: 9.10},
		{target: "a", value: 8.12},
		{target: "o", value: 7.68},
		{target: "i", value: 7.31},
		{target: "n", value: 6.95},
		{target: "s", value: 6.28},
		{target: "r", value: 6.02},
		{target: "d", value: 4.32},
		{target: "l", value: 3.98},
		{target: "u", value: 2.88},
		{target: "c", value: 2.71},
		{target: "m", value: 2.61},
		{target: "f", value: 2.30},
		{target: "y", value: 2.11},
		{target: "w", value: 2.09},
		{target: "g", value: 2.03},
		{target: "p", value: 1.82},
		{target: "b", value: 1.49},
		{target: "v", value: 1.11},
		{target: "k", value: 0.69},
		{target: "x", value: 0.17},
		{target: "q", value: 0.11},
		{target: "j", value: 0.10},
		{target: "z", value: 0.07},
	}

	freq := func(input, target string) float64 {
		count := 0
		size := 0

		for _, c := range strings.Split(input, "") {
			if c == " " {
				continue
			}

			size++
			if strings.ToLower(c) == target {
				count++
			}
		}

		return float64(count) / float64(size) * 100
	}

	score := float64(0)
	for _, e := range entries {
		score += math.Abs(e.value - freq(input, e.target))
	}

	return score
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type candidate struct {
	key   byte
	plain string
	score float64
}

func main() {
	fileContent, err := os.ReadFile("./4.txt")
	panicIfErr(err)

	candidates := make([]candidate, 0)
	for _, in := range strings.Split(string(fileContent), "\n") {
		candidates = append(candidates, decipher(in)...)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score < candidates[j].score
	})

	for _, c := range candidates[0:10] {
		fmt.Printf("`%s` %f\n", c.plain, c.score)
	}
}

func cipher(plain []byte, key []byte) []byte {
	keySize := len(key)
	buf := &bytes.Buffer{}

	for i, p := range plain {
		err := buf.WriteByte(p ^ key[i%keySize])
		panicIfErr(err)
	}

	return buf.Bytes()
}

func decipher(encrypted string) []candidate {
	candidates := make([]candidate, 0)
	for i := 0; i < int(math.Pow(2, 8)); i++ {
		key := byte(i)
		encrypted, err := hex.DecodeString(encrypted)
		panicIfErr(err)

		plainBuilder := &strings.Builder{}

		for _, b := range encrypted {
			plainBuilder.WriteByte(b ^ key)
		}

		plain := plainBuilder.String()
		if !utf8.ValidString(plain) {
			// filter out invalid utf-8 text
			continue
		}

		candidates = append(candidates, candidate{
			key:   key,
			plain: plain,
			score: frequencyScore(plain),
		})
	}

	return candidates
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
	expectedDistribution := map[string]float64{
		"e": 12.02,
		"t": 9.10,
		"a": 8.12,
		"o": 7.68,
		"i": 7.31,
		"n": 6.95,
		"s": 6.28,
		"r": 6.02,
		"d": 4.32,
		"l": 3.98,
		"u": 2.88,
		"c": 2.71,
		"m": 2.61,
		"f": 2.30,
		"y": 2.11,
		"w": 2.09,
		"g": 2.03,
		"p": 1.82,
		"b": 1.49,
		"v": 1.11,
		"k": 0.69,
		"x": 0.17,
		"q": 0.11,
		"j": 0.10,
		"z": 0.07,
	}

	distrib := func(input string) map[string]float64 {
		distribution := make(map[string]float64)
		size := 0
		for _, c := range strings.Split(input, "") {
			size++
			c = strings.ToLower(c)
			count, present := distribution[c]
			if present {
				distribution[c] = count + 1
			} else {
				distribution[c] = 1
			}

		}

		for char, count := range distribution {
			distribution[char] = count / float64(size) * 100
		}

		return distribution
	}

	chiSquare := func(expected, actual map[string]float64) (score float64) {
		score = 0

		for k, actualValue := range actual {
			if k == " " || k == "\n" {
				// we don't have blank space in our expected distribution, so skipping this for evaluation
				continue
			}

			expectedValue, ok := expected[k]
			if !ok {
				// artificially make this unexpected value very rare
				expectedValue = 0.001
			}

			score += math.Pow(actualValue-expectedValue, 2) / expectedValue
		}

		return score
	}

	return chiSquare(expectedDistribution, distrib(input))
}

func hammingDistance(lhs, rhs []byte) (distance int) {
	if len(lhs) != len(rhs) {
		//  let's make this assumption
		panic("trying to measure hamming distance of different length data")
	}

	distance = 0
	for i, l := range lhs {
		r := rhs[i]

		for i := 0; i < 8; i++ {
			n := byte(1 << i)
			if l&n != r&n {
				distance++
			}
		}
	}

	return distance
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

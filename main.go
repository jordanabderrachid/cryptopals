package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"
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
			if k == " " {
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

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexToBase64(t *testing.T) {
	in, expectedOutput := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d", "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	output := hexToBase64(in)

	assert.Equal(t, expectedOutput, output)
}

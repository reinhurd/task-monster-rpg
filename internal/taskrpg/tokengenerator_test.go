package taskrpg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringBytesMaskImpr(t *testing.T) {
	length := 10
	result := RandStringBytesMaskImpr(length)

	assert.Equal(t, length, len(result), "The length of the result should match the input length")

	for _, letter := range result {
		assert.True(t, strings.Contains(letterBytes, string(letter)), "The result should only contain letters from the defined set")
	}

	// testing randomness by generating 1000 strings and checking for any duplicates
	uniqueStrings := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		s := RandStringBytesMaskImpr(length)
		uniqueStrings[s] = true
	}
	assert.True(t, len(uniqueStrings) > 1, "Random function should generate unique strings")
}

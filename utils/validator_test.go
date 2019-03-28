package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAddress(t *testing.T) {
	assert.EqualValues(t, false, IsAddress("32d18e882ca24ba987662fb0de9052d09a29e5af"))
	assert.EqualValues(t, false, IsAddress("0x32d18e882ca24ba987662fb0de9052d09a29e5az"))
	assert.EqualValues(t, true, IsAddress("0x32d18e882ca24ba987662fb0de9052d09a29e5af"))
}

func TestIsTransaction(t *testing.T) {
	assert.EqualValues(t, false, IsTransaction("8d1ef09904f5b4c033b215848fb0c8de55bc435fb5bab22d1167eb34b3ea3ee7"))
	assert.EqualValues(t, false, IsTransaction("0x8d1ef09904f5b4c033b215848fb0c8de55bc435fb5bab22d1167eb34b3ea3eez"))
	assert.EqualValues(t, true, IsTransaction("0x8d1ef09904f5b4c033b215848fb0c8de55bc435fb5bab22d1167eb34b3ea3ee7"))
}

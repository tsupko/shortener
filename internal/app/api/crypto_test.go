package api

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeID(t *testing.T) {
	id, _ := decodeID("048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21")
	assert.Equal(t, "76543210", strconv.Itoa(int(id)))
}

func TestDecodeID2(t *testing.T) {
	id, _ := decodeID("d433b6ff37b4af30f2c6ea729330dfa41bc2be15ea0a00981f7d94a6086adc2a7b30bc91")
	assert.Equal(t, "3560158975", strconv.Itoa(int(id)))
}

func TestDecodeID3(t *testing.T) {
	id, _ := decodeID("d1e40328018721d2670537e4a7dccd8b03200ea9b5d124d56b1920819b1ce61480b96aea")
	assert.Equal(t, "3521381160", strconv.Itoa(int(id)))
}

func TestEncodeID(t *testing.T) {
	id := encodeID(uint32(76543210))
	assert.Equal(t, "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21", id)
}

func TestGetSignedID(t *testing.T) {
	id := getSignedUserID()
	assert.NotEqual(t, getSignedUserID(), id)
	assert.Equal(t, 72, len(id))
}

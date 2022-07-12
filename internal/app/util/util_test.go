package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert.NotSame(t, GenerateUniqueID(), GenerateUniqueID())
	assert.Len(t, GenerateUniqueID(), shortURLLength)
}

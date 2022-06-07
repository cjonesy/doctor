package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanHome(t *testing.T) {
	result, err := CleanHome("~/foo/bar")

	assert.NoError(t, err)
	assert.NotContains(t, result, "~")
}

package config

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigRunBasic(t *testing.T) {
	// Create a test config
	f, err := os.CreateTemp("", ".doctor.yml")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	testCfg := fmt.Sprintf(`
checks:
  - description: This is a basic test
    fix: Test fix
    id: basic-test
    path: %s
    type: file-exists
`, f.Name())

	os.WriteFile(f.Name(), []byte(testCfg), 0644)

	testConfig := Config{
		Path:    f.Name(),
		Verbose: true,
	}

	result := testConfig.Run()

	assert.NoError(t, result)
}

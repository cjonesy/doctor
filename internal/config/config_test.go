package config

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigRunBasic(t *testing.T) {
	// Create a test config
	f, err := ioutil.TempFile("", ".doctor.yml")
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

	ioutil.WriteFile(f.Name(), []byte(testCfg), 0644)

	testConfig := Config{
		Path:    f.Name(),
		Verbose: true,
	}

	result := testConfig.Run()

	assert.NoError(t, result)
}

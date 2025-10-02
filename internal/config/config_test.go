package config

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/cjonesy/doctor/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestConfigRunBasic(t *testing.T) {
	ctx := context.Background()
	log := logger.New(false, false)

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
		Logger:  log,
	}

	result := testConfig.Run(ctx)

	assert.NoError(t, result)
}

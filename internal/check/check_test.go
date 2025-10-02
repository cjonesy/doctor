package check

import (
	"context"
	"os"
	"syscall"
	"testing"

	"github.com/cjonesy/doctor/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestCheckRunBasic(t *testing.T) {
	ctx := context.Background()
	log := logger.New(false, false)

	// Create a test file
	f, err := os.CreateTemp("", "test.file")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	os.WriteFile(f.Name(), []byte(`foo=bar`), 0644)

	testCheckConfig := Check{
		Description: "This is a basic test",
		Type:        "file-exists",
		Path:        f.Name(),
		Verbose:     true,
		Logger:      log,
	}

	testRun := testCheckConfig.Run(ctx)

	assert.NoError(t, testRun)
}

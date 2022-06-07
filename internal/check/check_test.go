package check

import (
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRunBasic(t *testing.T) {
	// Create a test file
	f, err := ioutil.TempFile("", "test.file")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte(`foo=bar`), 0644)

	testCheckConfig := Check{
		Description: "This is a basic test",
		Type:        "file-exists",
		Path:        f.Name(),
		Verbose:     true,
	}

	testRun := testCheckConfig.Run()

	assert.NoError(t, testRun)
}

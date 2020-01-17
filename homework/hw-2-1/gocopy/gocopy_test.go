package gocopy

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const autogenSize = 1 * 1024 * 1024 // 1 MByte

var copyTests = []struct {
	outputFileName   string
	limit            int
	offset           int
	expectedFileSize int64
}{
	{
		"test-copy-1.dat",
		-1,
		0,
		autogenSize,
	},
}

func deleteFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("Error while deleting file %s: %s", path, err.Error())
	}
}

func TestCopy(t *testing.T) {
	var err error

	// Create test file with random content
	autogen, err := os.Create("test-autogen.dat")
	if err != nil {
		log.Fatalf("Could not open file for writing: %s", err.Error())
	}
	random, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatalf("Could not open random content generator: %s", err.Error())
	}
	written, err := io.CopyN(autogen, random, autogenSize)
	if err != nil {
		log.Fatalf("Could not create random content: %s", err.Error())
	}
	if written != autogenSize {
		log.Fatalf("Number of written bytes is less than expected: %d", written)
	}
	autogen.Close()
	random.Close()

	for _, test := range copyTests {
		err = Copy("test-autogen.dat", test.outputFileName, test.limit, test.offset)
		assert.Nil(t, err, "There should be no errors in this test")

		// Delete copy file
		deleteFile(test.outputFileName)
	}

	// Delete test file with content
	deleteFile("test-autogen.dat")
}

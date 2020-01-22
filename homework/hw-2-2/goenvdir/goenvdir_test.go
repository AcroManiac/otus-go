package goenvdir

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var envs = []struct {
	dirName  string
	fileName string
	content  string
}{
	{
		"testdir",
		"CITY",
		"Shi'Kahr",
	},
	{
		"testdir",
		"PLANET",
		"Vulcan",
	},
	{
		"testdir",
		"USER",
		"Spock",
	},
}

func createTestDir() error {
	if len(envs) == 0 {
		return errors.New("tests are empty")
	}

	if err := os.Mkdir(envs[0].dirName, 0755); err != nil {
		return err
	}

	for _, e := range envs {
		file, err := os.Create(filepath.Join(e.dirName, e.fileName))
		if err != nil {
			return err
		}
		_, err = file.WriteString(e.content)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func TestReadDir(t *testing.T) {
	if err := createTestDir(); err != nil {
		t.Error(err)
	}

	dirName := envs[0].dirName
	assert.DirExists(t, envs[0].dirName, "Directory should exist")

	// List files in test directory
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		t.Error(err)
	}

	// Check if file exists
	for _, file := range files {
		assert.FileExists(t, filepath.Join(dirName, file.Name()), "File should exist")

		// Check file content
		content, err := ioutil.ReadFile(filepath.Join(dirName, file.Name()))
		if err != nil {
			t.Error(err)
		}

		// Find expected value
		var expected string
		for _, e := range envs {
			if e.fileName == file.Name() {
				expected = e.content
				break
			}
		}
		assert.Equal(t, expected, string(content), "Strings should be equal")
	}

	// Remove test directory
	if err := os.RemoveAll(dirName); err != nil {
		t.Error(err)
	}
}

func TestRunCmd(t *testing.T) {
	//
}

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

var commands = []struct {
	command  []string
	output   []string
	exitCode int
}{
	{
		[]string{"printenv", "USER"},
		[]string{"Spock"},
		0,
	},
	{
		[]string{"printenv", "CITY", "PLANET", "USER"},
		[]string{"Shi'Kahr", "Vulcan", "Spock"},
		0,
	},
}

func createTestDir() error {
	if len(envs) == 0 {
		return errors.New("tests are empty")
	}

	dirName := envs[0].dirName
	if err := os.Mkdir(dirName, 0755); err != nil {
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

	// List files in test directory
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	// Check if file exists
	for _, file := range files {
		// Check file content
		content, err := ioutil.ReadFile(filepath.Join(dirName, file.Name()))
		if err != nil {
			return err
		}

		// Find expected value
		var expected string
		for _, e := range envs {
			if e.fileName == file.Name() {
				expected = e.content
				break
			}
		}
		if expected != string(content) {
			return errors.New("error in file content")
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

	env, err := ReadDir(dirName)
	if err != nil {
		t.Error(err)
	}

	// Check environment
	for _, e := range envs {
		value, ok := env[e.fileName]
		if !ok {
			t.Errorf("no such environment variable: %s", e.fileName)
		}
		if value != e.content {
			t.Errorf("error in environment variable content: %s", e.content)
		}
	}

	// Remove test directory
	if err := os.RemoveAll(dirName); err != nil {
		t.Error(err)
	}
}

func TestRunCmd(t *testing.T) {
	if err := createTestDir(); err != nil {
		t.Error(err)
	}

	dirName := envs[0].dirName
	assert.DirExists(t, dirName, "Directory should exist")

	env, err := ReadDir(dirName)
	if err != nil {
		t.Errorf("Error reading environment directory: %s", err.Error())
	}

	for _, command := range commands {
		code := RunCmd(command.command, env)
		assert.Equal(t, command.exitCode, code, "Command exit code should be equal to expected")
	}

	// Remove test directory
	if err := os.RemoveAll(dirName); err != nil {
		t.Error(err)
	}
}

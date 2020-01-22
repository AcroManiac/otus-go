package goenvdir

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ReadDir reads files in a directory with path dir and fills
// an environment map with key = filename and value = file content.
// In case of file operation errors ReadDir returns them to calling function
func ReadDir(dir string) (map[string]string, error) {
	env := make(map[string]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, file := range files {
		key := file.Name()
		value, err := ioutil.ReadFile(filepath.Join(dir, key))
		if err != nil {
			return env, err
		}

		env[key] = string(value)
	}

	return env, nil
}

// RunCmd runs external command with arguments from cmd slice and
// environment variables defined in env map. It returns the exit code
// of external command
func RunCmd(cmd []string, env map[string]string) int {
	args := cmd[1:]
	command := exec.Command(cmd[0], args...)

	// Make map from slice
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envMap[pair[0]] = pair[1]
	}

	// Redefine environment with env values
	for key, value := range env {
		_, ok := envMap[key]
		if ok {
			// Remove variable if it is empty
			// See http://manpages.ubuntu.com/manpages/bionic/man8/envdir.8.html
			if value == "" {
				delete(envMap, key)
				continue
			}
		}
		envMap[key] = value
	}

	// Make slice from map
	newEnv := make([]string, 0, len(envMap))
	for key, value := range envMap {
		newEnv = append(newEnv, strings.Join([]string{key, value}, "="))
	}

	// Set new environment and standard streams
	command.Env = newEnv
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	// Unfortunately, there is no platform independent way to get the exit code
	// in the error case. That's also the reason why it isn't part of the API.
	// The following snippet will work with Linux, but I haven't tested it on other platforms
	if err := command.Start(); err != nil {
		log.Printf("Error starting command : %v", err)
		return 1
	}

	var code int
	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
				log.Printf("Exit status: %d", code)
			}
		} else {
			log.Printf("Error waiting command: %v", err)
			code = 1
		}
	}

	return code
}

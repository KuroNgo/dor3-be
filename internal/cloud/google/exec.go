package google

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteTask(execCmd string) error {
	// execute command string parts, delimited by space
	execParts := strings.Fields(execCmd)
	if len(execParts) == 0 {
		return errors.New("empty command")
	}

	// executable name
	execName := execParts[0]

	// execute command parameters
	execParams := execParts[1:]

	// execute command instance
	cmd := exec.Command(execName, execParams...)

	// run execute command instance
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}

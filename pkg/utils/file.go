package utils

import "os/exec"

func FindGoModPath() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

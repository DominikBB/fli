package history

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Fetcher interface {
	GetLastCommand() (string, error)
}

type History struct {
	Shell string
}

var _ Fetcher = (*History)(nil)

const (
	ShellNu = "nu"
)

var (
	ErrUnsupportedShell    = errors.New("unsupported shell")
	ErrHistoryPathNotFound = errors.New("unable to locate history file")
)

func (h History) GetLastCommand() (string, error) {
	switch h.Shell {
	case ShellNu:
		historyPath, err := getNushellHistoryPath()
		if err != nil {
			return "", ErrHistoryPathNotFound
		}
		cmd, err := readLastLine(filepath.Join(historyPath))
		if err != nil {
			return "", ErrHistoryPathNotFound
		}
		return cmd, nil
	default:
		return "", ErrUnsupportedShell
	}
}

func getNushellHistoryPath() (string, error) {
	cmd := exec.Command("nu", "-c", "echo $nu.default-config-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Nushell config directory: %w", err)
	}

	configDir := strings.TrimSpace(string(output))
	historyPath := filepath.Join(configDir, "history.txt")
	return historyPath, nil
}

// readLastLine reads the last non-empty line of a file.
// Returns an error if the file can't be opened or is empty.
func readLastLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			lastLine = line
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if lastLine == "" {
		return "", errors.New("no non-empty lines found")
	}

	return lastLine, nil
}

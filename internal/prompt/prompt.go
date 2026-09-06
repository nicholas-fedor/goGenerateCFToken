// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package prompt provides shared user-input helpers for interactive CLI prompts.
package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrEmptyInput indicates the user submitted an empty response.
var ErrEmptyInput = errors.New("input cannot be empty")

// Prompt displays a message to stderr and reads a line from stdin.
//
// Parameters:
//   - message: The prompt text to display.
//
// Returns:
//   - string: The user's input (trimmed of trailing newline).
//   - error: Non-nil if reading from stdin fails or input is empty.
func Prompt(message string) (string, error) {
	fmt.Fprint(os.Stderr, message)

	result, err := readLine()
	if err != nil {
		return "", err
	}

	if result == "" {
		return "", ErrEmptyInput
	}

	return result, nil
}

// Secret displays a message to stderr and reads a secret from stdin.
// On a TTY the input is not echoed. Piped stdin is read as a single line.
//
// Parameters:
//   - message: The prompt text to display.
//
// Returns:
//   - string: The secret (trimmed of surrounding whitespace).
//   - error: Non-nil if reading fails or the secret is empty.
func Secret(message string) (string, error) {
	fmt.Fprint(os.Stderr, message)

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		secret, err := term.ReadPassword(fd)

		fmt.Fprintln(os.Stderr)

		if err != nil {
			return "", fmt.Errorf("read secret: %w", err)
		}

		result := strings.TrimSpace(string(secret))
		if result == "" {
			return "", ErrEmptyInput
		}

		return result, nil
	}

	result, err := readLine()
	if err != nil {
		return "", err
	}

	if result == "" {
		return "", ErrEmptyInput
	}

	return result, nil
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	return strings.TrimSpace(input), nil
}

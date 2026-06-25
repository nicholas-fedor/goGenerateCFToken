// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package logging provides structured logging via zerolog and plain-text
// user-facing output helpers.
package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Setup initializes zerolog with appropriate level and formatting.
//
// In debug mode, output is structured JSON to stderr with caller info
// and error stack traces enabled. In normal mode, zerolog output goes
// to stderr and user-facing messages go to stdout via Printf/Println.
//
// Parameters:
//   - isDebug: If true, enables debug-level logging with structured JSON output.
func Setup(isDebug bool) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign

	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		log.Logger = log.Output(os.Stderr).With().Caller().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		log.Logger = log.Output(os.Stderr)
	}
}

// Logger returns the global zerolog logger.
// Use for structured logging (Debug, Error with Stack).
// In debug mode this produces JSON to stderr.
//
// Returns:
//   - zerolog.Logger: The global logger instance.
func Logger() zerolog.Logger {
	return log.Logger
}

// SetLevel updates the global log level at runtime.
//
// Parameters:
//   - level: The new log level to set.
func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// Printf writes a plain text message to stdout with a trailing newline.
// Use for user-facing output in normal (non-debug) mode.
//
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
func Printf(format string, v ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format+"\n", v...)
}

// Println writes a plain text message to stdout with a trailing newline.
// Use for user-facing output in normal (non-debug) mode.
//
// Parameters:
//   - v: Variadic arguments to print.
func Println(v ...any) {
	_, _ = fmt.Fprintln(os.Stdout, v...)
}

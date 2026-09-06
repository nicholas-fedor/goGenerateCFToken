// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package logging provides structured logging via zerolog and plain-text
// user-facing output helpers.
package logging

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var quiet atomic.Bool

// Setup initializes zerolog at the requested level.
//
// Debug (and trace) use structured JSON to stderr with caller info.
// Other levels write zerolog output to stderr without caller info.
// User-facing messages go to stdout via Printf/Println unless quiet.
//
// Parameters:
//   - level: Global zerolog level to apply.
func Setup(level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign

	zerolog.SetGlobalLevel(level)

	if level <= zerolog.DebugLevel {
		log.Logger = log.Output(os.Stderr).With().Caller().Logger()
	} else {
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

// SetQuiet suppresses Printf/Println user-facing output when true.
//
// Parameters:
//   - q: Whether incidental user-facing output should be discarded.
func SetQuiet(q bool) {
	quiet.Store(q)
}

// Printf writes a plain text message to stdout with a trailing newline.
// Use for incidental user-facing output. Honors SetQuiet.
//
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
func Printf(format string, v ...any) {
	if quiet.Load() {
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, format+"\n", v...)
}

// Println writes a plain text message to stdout with a trailing newline.
// Use for incidental user-facing output. Honors SetQuiet.
//
// Parameters:
//   - v: Variadic arguments to print.
func Println(v ...any) {
	if quiet.Load() {
		return
	}

	_, _ = fmt.Fprintln(os.Stdout, v...)
}

// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package logging

import "fmt"

// redactionMarker is what a Secret renders as in every output path.
const redactionMarker = "[REDACTED]"

// Secret wraps a sensitive value so that it cannot be printed by accident.
//
// It is a guard rail, not a container for credentials the program needs to read back: a
// Secret deliberately offers no way to reveal the value again, so it is only useful for
// values on their way to a log call. The value it protects is never rendered, through
// fmt formatting, text or JSON marshalling, or the zerolog Stringer and Interface hooks.
//
// A redacting zerolog hook is not an alternative here. A hook runs when the event is
// written, after the fields have already been serialised, and zerolog exposes no way to
// read a field back or rewrite it. A hook can only append, which leaves the original value
// in the output while adding a second key that looks redacted, so it is worse than no
// protection at all. Removing what was already written requires a writer that parses each
// line, which is deliberately not done here.
//
// The current call sites in the credential path do not need this type: they log names and
// paths rather than values. It exists so that future code logging a credential by mistake
// has a safe way to do so.
//
// One verb is out of reach. fmt handles %p before it looks for any method, so a value type
// always renders its fields inside fmt's own "%!p(...)" diagnostic and no method on this
// type can intervene. %p addresses pointers, so no realistic logging call would use it on a
// Secret, but the limitation is recorded here rather than left to be discovered.
type Secret struct {
	value string
}

// NewSecret wraps a value so that logging it produces only the redaction marker.
//
// Parameters:
//   - value: The sensitive value to protect.
//
// Returns:
//   - Secret: A value that renders as the redaction marker in every output path.
func NewSecret(value string) Secret {
	return Secret{value: value}
}

// String returns the redaction marker rather than the value, so a Secret printed through
// fmt with %s, %v or %q, or through zerolog's Stringer hook, cannot leak.
//
// Returns:
//   - string: The redaction marker.
func (s Secret) String() string {
	return redactionMarker
}

// GoString returns the redaction marker, so %#v formatting is also safe even for a caller
// that reaches for a reflection-based pretty printer rather than fmt.
//
// Returns:
//   - string: The redaction marker.
func (s Secret) GoString() string {
	return redactionMarker
}

// Format implements fmt.Formatter, which is what makes the type safe for every verb rather
// than only the string-like ones. fmt consults Stringer for %s, %v, %q and the hex verbs,
// but a mismatched verb such as %d falls through to reflection on the wrapped field and
// prints the credential inside fmt's own "%!d(string=...)" diagnostic. Handling every verb
// here closes that route.
//
// Parameters:
//   - f: The formatter state to write the marker to.
//   - c: The verb, accepted for interface conformance and otherwise ignored.
func (s Secret) Format(f fmt.State, _ rune) {
	_, _ = f.Write([]byte(redactionMarker))
}

// MarshalText returns the redaction marker, so encoding through
// encoding.TextMarshaler is safe.
//
// Returns:
//   - []byte: The redaction marker.
//   - error: Always nil.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(redactionMarker), nil
}

// MarshalJSON returns the redaction marker as a JSON string, so marshalling a struct that
// embeds a Secret cannot leak.
//
// Returns:
//   - []byte: The marker as a quoted JSON string.
//   - error: Always nil.
func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte(`"` + redactionMarker + `"`), nil
}

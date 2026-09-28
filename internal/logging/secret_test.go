// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package logging

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecret_RendersMarkerThroughEveryOutputPath(t *testing.T) {
	const raw = "super-secret-value"

	secret := NewSecret(raw)

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, redactionMarker, secret.String())
	})

	t.Run("GoString", func(t *testing.T) {
		assert.Equal(t, redactionMarker, secret.GoString())
	})

	t.Run("MarshalText", func(t *testing.T) {
		got, err := secret.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, redactionMarker, string(got))
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		got, err := secret.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"`+redactionMarker+`"`, string(got))
	})
}

func TestSecret_NeverRendersTheValue(t *testing.T) {
	const raw = "super-secret-value"

	secret := NewSecret(raw)

	formats := []string{
		"%s", "%v", "%q", "%x", "%X", "%#v", "%+v", "%T",
		// Verbs that do not match a string would otherwise fall through to reflection on
		// the wrapped field and print the credential in fmt's own diagnostic, so they are
		// covered deliberately.
		"%d", "%f", "%t", "%c", "%U",
	}

	for _, format := range formats {
		t.Run(fmt.Sprintf("fmt.Sprintf(%q)", format), func(t *testing.T) {
			assert.NotContains(t, fmt.Sprintf(format, secret), raw)
		})
	}
}

func TestSecret_PercentPUsesFmtDiagnostic(t *testing.T) {
	const raw = "super-secret-value"

	// This is a characterisation test, not a safety guarantee, and it is the inverse of the
	// one above: the assertion is that the credential DOES appear. fmt handles the p verb
	// before it looks for any method, so a value type always renders its fields inside its
	// own diagnostic and no method on Secret can intervene.
	//
	// It exists so that a change in fmt behaviour fails here and prompts the limitation
	// recorded on the type to be revisited, rather than the hole being discovered later.
	got := fmt.Sprintf("%p", NewSecret(raw))

	assert.Contains(t, got, raw, "if this no longer leaks, fmt behaviour changed and the type's documentation should be updated")
	assert.Contains(t, got, "%!p", "the malformed diagnostic is what makes the leak visible to a reader")
}

func TestSecret_ProtectsThePlainTextHelpers(t *testing.T) {
	const raw = "super-secret-value"

	// logging.Printf and Println take ...any and format with fmt, so a Secret passed to
	// either renders as the marker rather than the credential. This is the path a future
	// Printf("resolved key %s", key) mistake would take.
	assert.NotContains(t, fmt.Sprintf("%s", NewSecret(raw)), raw)
	assert.NotContains(t, fmt.Sprint(NewSecret(raw)), raw)
}

func TestSecret_ProtectsEmbeddedStructs(t *testing.T) {
	const raw = "super-secret-value"

	type record struct {
		Name  string `json:"name"`
		Token Secret `json:"token"`
	}

	got, err := json.Marshal(record{Name: "svc", Token: NewSecret(raw)})
	require.NoError(t, err)
	assert.NotContains(t, string(got), raw)
	assert.Contains(t, string(got), redactionMarker)
}

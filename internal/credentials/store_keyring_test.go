// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keyring "github.com/zalando/go-keyring"
)

func TestNewKeyringStore_DoesNotMutateLiveCredential(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, keyring.Set(KeyringService, KeyringUser, "real-token"))

	store := NewKeyringStore(KeyringService, KeyringUser)

	got, err := keyring.Get(KeyringService, KeyringUser)
	require.NoError(t, err)
	assert.Equal(t, "real-token", got)
	assert.True(t, store.Available())
}

func TestNewKeyringStore_UnsupportedBackend(t *testing.T) {
	keyring.MockInitWithError(keyring.ErrUnsupportedPlatform)
	t.Cleanup(keyring.MockInit)

	store := NewKeyringStore(KeyringService, KeyringUser)

	assert.False(t, store.Available())
}

func TestKeyringStore_SetGetDelete(t *testing.T) {
	keyring.MockInit()

	store := NewKeyringStore("test-service", "test-user")
	require.True(t, store.Available())

	require.NoError(t, store.Set("stored-key"))

	got, err := store.Get()
	require.NoError(t, err)
	assert.Equal(t, "stored-key", got)

	require.NoError(t, store.Delete())

	_, err = store.Get()
	require.Error(t, err)
	require.ErrorIs(t, err, keyring.ErrNotFound)

	require.NoError(t, store.Delete())
}

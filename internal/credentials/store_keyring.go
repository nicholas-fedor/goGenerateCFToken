// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	keyring "github.com/zalando/go-keyring"
)

// KeyringStore provides low-level OS keyring operations for a service/user pair.
type KeyringStore struct {
	service   string
	user      string
	available bool
}

// probeUser is a dedicated keyring user used only to detect backend
// availability. It is never the live credential slot.
const probeUser = "__gogeneratecftoken_probe__"

// NewKeyringStore creates a KeyringStore for the given service and user.
// Availability is probed with a non-mutating Get against a dummy user so the
// live credential is never overwritten or deleted.
//
// Parameters:
//   - service: The keyring service name.
//   - user: The keyring user name.
//
// Returns:
//   - *KeyringStore: A new keyring store instance.
func NewKeyringStore(service, user string) *KeyringStore {
	return &KeyringStore{
		service:   service,
		user:      user,
		available: probeAvailable(service),
	}
}

// probeAvailable reports whether a keyring backend can be reached.
//
// Parameters:
//   - service: The keyring service name used for the probe Get.
//
// Returns:
//   - bool: True when the backend responds (including "not found").
func probeAvailable(service string) bool {
	_, err := keyring.Get(service, probeUser)

	return err == nil || errors.Is(err, keyring.ErrNotFound)
}

// Get retrieves the stored key from the OS keyring.
//
// Returns:
//   - string: The stored key.
//   - error: Non-nil if the key is not found or keyring access fails.
func (s *KeyringStore) Get() (string, error) {
	key, err := keyring.Get(s.service, s.user)
	if err != nil {
		log.Debug().Err(err).Msg("failed to read API key from OS keyring")

		return "", fmt.Errorf("keyring get failed: %w", err)
	}

	log.Debug().Msg("resolved API key from OS keyring")

	return key, nil
}

// Set stores the key in the OS keyring.
//
// Parameters:
//   - key: The key to store.
//
// Returns:
//   - error: Non-nil if keyring storage fails.
func (s *KeyringStore) Set(key string) error {
	log.Debug().Msg("storing API key in OS keyring")

	err := keyring.Set(s.service, s.user, key)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to store API key in keyring")

		return fmt.Errorf("keyring set failed: %w", err)
	}

	return nil
}

// Delete removes the key from the OS keyring. It returns nil if the key did
// not exist (idempotent).
//
// Returns:
//   - error: Non-nil if keyring deletion fails for reasons other than not found.
func (s *KeyringStore) Delete() error {
	log.Debug().Msg("removing API key from OS keyring")

	err := keyring.Delete(s.service, s.user)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		log.Error().Stack().Err(err).Msg("failed to remove API key from keyring")

		return fmt.Errorf("keyring delete failed: %w", err)
	}

	if errors.Is(err, keyring.ErrNotFound) {
		log.Debug().Msg("API key not found in keyring - nothing to remove")
	}

	return nil
}

// Available reports whether the keyring backend was reachable at construction
// time.
//
// Returns:
//   - bool: True when a keyring backend was available.
func (s *KeyringStore) Available() bool {
	return s.available
}

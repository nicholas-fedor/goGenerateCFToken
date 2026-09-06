// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentials

import "context"

type mockProvider struct {
	key string
	err error
}

func (m mockProvider) Resolve(_ context.Context) (string, error) {
	return m.key, m.err
}

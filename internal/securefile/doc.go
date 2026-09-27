// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package securefile writes files that hold secrets.
//
// A secret on disk needs two properties that os.WriteFile does not provide. The mode
// argument of os.WriteFile is only applied when the kernel creates the file, so writing
// over a file that already exists silently keeps whatever mode it had, which leaves a
// secret at whatever mode a backup, a copy or an earlier mistake left behind. And
// os.WriteFile truncates the destination before writing, so an interrupted write leaves an
// empty or partial file where the secret should be.
//
// WriteFile here writes to a temporary file in the same directory and renames it over the
// destination, which also means a symlink at the destination is replaced rather than
// written through. See WriteFile for how far the replacement is atomic on each platform.
package securefile

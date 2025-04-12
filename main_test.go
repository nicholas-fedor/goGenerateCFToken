/*
Copyright © 2025 Nicholas Fedor <nick@nickfedor.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"os"
	"testing"

	"github.com/nicholas-fedor/goGenerateCFToken/cmd"
)

func TestMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	oldExit := osExit
	osExit = func(code int) {
		if code != 0 {
			t.Errorf("main() tried to exit with code %d", code)
		}
	}

	defer func() { osExit = oldExit }()

	main()
}

var osExit = os.Exit

func init() {
	cmdExecute = func() {
	}
}

var cmdExecute = cmd.Execute

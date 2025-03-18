package main

import (
	"os"
	"testing"

	"github.com/nicholas-fedor/goGenerateCFToken/cmd"
)

func TestMain(t *testing.T) {
	// Since main() only calls cmd.Execute(), we test that it doesn't panic
	// and rely on cmd package tests for detailed behavior.
	// This is a minimal test as main() has no direct logic to test.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	// Replace os.Exit to prevent test termination
	oldExit := osExit
	osExit = func(code int) {
		if code != 0 {
			t.Errorf("main() tried to exit with code %d", code)
		}
	}

	defer func() { osExit = oldExit }()

	main()
}

// osExit is a variable to mock os.Exit.
var osExit = os.Exit

// Mock cmd.Execute to avoid actual execution in main_test.
func init() {
	cmdExecute = func() {
		// Do nothing, just return to allow testing main()
	}
}

var cmdExecute = cmd.Execute

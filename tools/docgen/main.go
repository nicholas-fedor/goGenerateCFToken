// Copyright (c) Nicholas Fedor 2026 <nick@nickfedor.com>
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:generate go run github.com/nicholas-fedor/gogeneratecftoken/v2/tools/docgen -out ../../www/content/cli-reference

// Command docgen generates Hugo-compatible Markdown documentation for the CLI command tree.
//
// Usage:
//
//	go run ./tools/docgen -out ../../www/content/cli-reference
package main

import (
	"flag"
	"log"

	"github.com/nicholas-fedor/gogeneratecftoken/v2/cmd"
)

func main() {
	out := flag.String("out", "./www/content/cli-reference", "Output directory")

	flag.Parse()

	generator, err := NewDocGenerator()
	if err != nil {
		log.Fatalf("create doc generator: %v", err)
	}

	root := cmd.Root()

	err = generator.Generate(root, *out)
	if err != nil {
		log.Fatal(err)
	}
}

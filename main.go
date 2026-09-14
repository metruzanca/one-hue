// Command one-hue generates static theme files for code editors from a
// shared color model. Run `one-hue --help` for usage.
package main

import (
	"os"

	"github.com/metruzanca/one-hue-theme/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
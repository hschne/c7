package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "c7",
	Short: "Context7 CLI — fetch up-to-date library documentation",
	Long: `c7 is a lightweight CLI for Context7 (context7.com).
Fetch current library documentation from the terminal.
No Node.js required, single binary, instant startup.`,
	Version: version,
}

func Execute(v string) {
	version = v
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

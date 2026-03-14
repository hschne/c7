package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "c7",
	Short: "Context7 CLI — fetch up-to-date library documentation",
	Long: `c7 is a lightweight CLI for Context7 (context7.com).
Fetch current library documentation from the terminal.
No Node.js required, single binary, instant startup.`,
}

func Execute(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

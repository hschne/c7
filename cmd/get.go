package cmd

import (
	"fmt"

	"github.com/hschne/c7/internal"
	"github.com/spf13/cobra"
)

var getTokens string

var getCmd = &cobra.Command{
	Use:   "get <library-name> <query>",
	Short: "Resolve a library and fetch docs in one step",
	Long: `One-shot: resolve library name then fetch docs. Easiest to use.
Caches the resolved library ID for faster repeat lookups.

Examples:
  c7 get rails "active record scopes"
  c7 get hotwire "form submission with turbo"
  c7 get kamal "deploy with secrets" --tokens 8000`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		libName := args[0]
		query := args[1]

		var libID string

		if entry, ok := internal.CacheLookup(libName); ok {
			libID = entry.ID
		} else {
			client := internal.NewClient()
			libs, err := client.Search(libName, query)
			if err != nil {
				return err
			}
			if len(libs) == 0 {
				return fmt.Errorf("no library found for: %s", libName)
			}
			best := libs[0]
			libID = best.ID
			internal.CacheSave(libName, best.ID, best.Name)
		}

		client := internal.NewClient()
		body, err := client.FetchDocs(libID, query, getTokens, "1", "")
		if err != nil {
			return err
		}

		internal.PrintDocs(body)
		return nil
	},
}

func init() {
	getCmd.Flags().StringVar(&getTokens, "tokens", "5000", "max tokens to return")
	rootCmd.AddCommand(getCmd)
}

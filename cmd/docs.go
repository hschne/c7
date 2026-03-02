package cmd

import (
	"github.com/hschne/c7/internal"
	"github.com/spf13/cobra"
)

var docsTokens string
var docsPage string
var docsTopic string

var docsCmd = &cobra.Command{
	Use:   "docs <library-id> <query>",
	Short: "Fetch docs for a known library ID",
	Long: `Fetch documentation for a known library ID.

Examples:
  c7 docs /rails/rails "how to use scopes"
  c7 docs /vercel/next.js "middleware" --topic routing --page 2
  c7 docs /hotwire-dev/turbo "streams" --tokens 10000`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		libraryID := args[0]
		query := args[1]

		client := internal.NewClient()
		body, err := client.FetchDocs(libraryID, query, docsTokens, docsPage, docsTopic)
		if err != nil {
			return err
		}

		internal.PrintDocs(body)
		return nil
	},
}

func init() {
	docsCmd.Flags().StringVar(&docsTokens, "tokens", "5000", "max tokens to return")
	docsCmd.Flags().StringVar(&docsPage, "page", "1", "page number (1-10)")
	docsCmd.Flags().StringVar(&docsTopic, "topic", "", "focus docs on a specific topic")
	rootCmd.AddCommand(docsCmd)
}

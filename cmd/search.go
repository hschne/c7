package cmd

import (
	"fmt"
	"strings"

	"github.com/hschne/c7/internal"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <library-name> [query]",
	Short: "Search for libraries by name",
	Long: `Search for libraries by name. Optional query improves relevance ranking.

Examples:
  c7 search rails
  c7 search "ruby on rails" "active record"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		libraryName := args[0]
		query := libraryName
		if len(args) > 1 {
			query = strings.Join(args[1:], " ")
		}

		client := internal.NewClient()
		libs, err := client.Search(libraryName, query)
		if err != nil {
			return err
		}

		if len(libs) == 0 {
			fmt.Println("No libraries found.")
			return nil
		}

		internal.CacheSave(libraryName, libs[0].ID, libs[0].Name)

		fmt.Printf("%-30s %-8s %s\n", "ID", "TRUST", "NAME")
		fmt.Println(strings.Repeat("─", 70))
		for _, l := range libs {
			fmt.Printf("%-30s %-8.1f %s\n", l.ID, l.TrustScore, l.Name)
			if l.Description != "" {
				for _, line := range internal.WrapText(l.Description, 60) {
					fmt.Printf("%-30s          %s\n", "", line)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

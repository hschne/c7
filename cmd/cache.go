package cmd

import (
	"fmt"

	"github.com/hschne/c7/internal"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the library lookup cache",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all cached library lookups",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CacheClear(); err != nil {
			return err
		}
		fmt.Println("Cache cleared.")
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)
}

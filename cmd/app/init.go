package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version         = "1.0.0"
	author          = "Herwono W. Wijaya <herwonowr@vulncode.com>"
	databaseVersion = 1

	rootCmd = &cobra.Command{
		Use:     "slackhell",
		Version: version,
		Short:   "Slack Web Shell C2.",
		Long:    "Slack Web Shell Command & Control\nBy " + author,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Slackhell v%s \n", version)
		},
		DisableFlagParsing: true,
	}
)

// Execute ...
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate(`{{print "Slackhell "}}{{printf "v%s\n" .Version}}`)
}

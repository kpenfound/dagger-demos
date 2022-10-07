package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	buildkitAddr string
	repository   string
	site         string
	repoRef      string

	rootCmd = &cobra.Command{
		Use:   "kyle-ci",
		Short: "A ci tool for the kyle organization",
		Long:  `kyle-ci is a CI tool used for ci operations at kyle inc.`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&buildkitAddr, "buildkit", "", "The address of the buildkit host")
	rootCmd.AddCommand(
		netlifyCmd,
	)

	netlifyCmd.Flags().StringVarP(&repository, "repository", "r", "", "git repository")
	netlifyCmd.MarkFlagRequired("repository")
	netlifyCmd.Flags().StringVarP(&site, "site", "s", "", "netlify site")
	netlifyCmd.Flags().StringVarP(&repoRef, "ref", "t", "", "git repository ref")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

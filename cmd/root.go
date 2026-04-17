package cmd

import (
	"fmt"
	"os"

	"github.com/TDiblik/git-forge/pkg/git"
	"github.com/spf13/cobra"
)

var (
	cloneHash  string
	sign       bool
	customDate string
	dryRun     bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "git-forge",
	Short: "A tool to demonstrate Git identity vulnerabilities",
	Long:  `git-forge is a CLI tool designed for educational and penetration-testing purposes to demonstrate vulnerabilities in unsigned Git repositories.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cloneHash, "clone", "", "Extract name, email, and dates from the specified SHA-1 hash")
	rootCmd.PersistentFlags().BoolVar(&sign, "sign", false, "Overrides the default --no-sign. Triggers the GPG generation process.")
	rootCmd.PersistentFlags().StringVar(&customDate, "date", "", "Explicitly sets the date, overriding system time and --clone.")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Prints the raw git and gpg commands without running them.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enables detailed logging.")
}

func resolveIdentity(authorFlag string, vipFlag string) (*git.Identity, error) {
	var finalID *git.Identity

	count := 0
	if cloneHash != "" {
		count++
	}
	if authorFlag != "" {
		count++
	}
	if vipFlag != "" {
		count++
	}

	if count > 1 {
		return nil, fmt.Errorf("fatal: --clone, --author, and --vip are mutually exclusive")
	}

	if cloneHash != "" {
		id, err := git.ResolveFromHash(cloneHash)
		if err != nil {
			return nil, err
		}
		finalID = id
	} else if authorFlag != "" {
		id, err := git.ParseAuthor(authorFlag)
		if err != nil {
			return nil, err
		}
		finalID = id
	} else if vipFlag != "" {
		id, err := git.ResolveVIP(vipFlag)
		if err != nil {
			return nil, err
		}
		finalID = id
	}

	if finalID != nil && customDate != "" {
		finalID.Date = customDate
	} else if finalID == nil && customDate != "" {
		finalID = &git.Identity{Date: customDate}
	}

	return finalID, nil
}

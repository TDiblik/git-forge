package cmd

import (
	"fmt"
	"strings"

	"github.com/TDiblik/git-forge/pkg/git"
	"github.com/TDiblik/git-forge/pkg/gpg"
	"github.com/spf13/cobra"
)

var (
	message   string
	author    string
	vip       string
	typoSquat string
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Creates a new spoofed commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveIdentity(author, vip)
		if err != nil {
			return err
		}

		if typoSquat != "" {
			if id == nil {
				id = &git.Identity{}
			}
			id.Email = git.TypoSquat(typoSquat)
		}

		opts := git.CommandOptions{
			Identity: id,
			DryRun:   dryRun,
			Verbose:  verbose,
			NoSign:   !sign,
		}

		if sign {
			if id == nil {
				return fmt.Errorf("GPG signing requires an identity (provide via --author, --clone, or --vip)")
			}
			gpgMgr, err := gpg.NewManager()
			if err != nil {
				return err
			}
			defer gpgMgr.Cleanup()

			keyID, err := gpgMgr.GenerateKey(id.Name, id.Email, id.Date, dryRun)
			if err != nil {
				return err
			}
			opts.SigningKey = keyID
			opts.GnuPGHome = gpgMgr.TempDir
		}

		gitArgs := []string{"commit", "-m", message}
		if sign {
			gitArgs = append(gitArgs, "-S")
		}

		return git.RunGitCommand(gitArgs, opts)
	},
}

func init() {
	vips := strings.Join(git.GetVIPs(), ", ")
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (required)")
	commitCmd.MarkFlagRequired("message")
	commitCmd.Flags().StringVarP(&author, "author", "a", "", "Explicit identity (Name <email>)")
	commitCmd.Flags().StringVar(&vip, "vip", "", fmt.Sprintf("Load identity from profile. Available: %s", vips))
	commitCmd.Flags().StringVar(&typoSquat, "typo-squat", "", "Generates a slightly misspelled email based on the input")
	rootCmd.AddCommand(commitCmd)
}

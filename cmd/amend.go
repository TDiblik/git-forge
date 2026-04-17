package cmd

import (
	"fmt"
	"strings"

	"github.com/TDiblik/git-forge/pkg/git"
	"github.com/TDiblik/git-forge/pkg/gpg"
	"github.com/spf13/cobra"
)

var (
	amendAuthor string
	amendVip    string
)

var amendCmd = &cobra.Command{
	Use:   "amend",
	Short: "Quick manipulation of HEAD",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveIdentity(amendAuthor, amendVip)
		if err != nil {
			return err
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

			keyID, err := gpgMgr.GenerateKey(id.Name, id.Email, dryRun)
			if err != nil {
				return err
			}
			opts.SigningKey = keyID
			opts.GnuPGHome = gpgMgr.TempDir
		}

		gitArgs := []string{"commit", "--amend", "--no-edit"}
		if id != nil {
			gitArgs = append(gitArgs, "--reset-author")
		}
		if sign {
			gitArgs = append(gitArgs, "-S")
		}

		return git.RunGitCommand(gitArgs, opts)
	},
}

func init() {
	vips := strings.Join(git.GetVIPs(), ", ")
	amendCmd.Flags().StringVarP(&amendAuthor, "author", "a", "", "Explicit identity (Name <email>)")
	amendCmd.Flags().StringVar(&amendVip, "vip", "", fmt.Sprintf("Load identity from profile. Available: %s", vips))
	rootCmd.AddCommand(amendCmd)
}

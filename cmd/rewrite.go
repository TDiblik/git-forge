package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/TDiblik/git-forge/pkg/git"
	"github.com/spf13/cobra"
)

var (
	rewriteAuthor string
	rewriteVip    string
)

var rewriteCmd = &cobra.Command{
	Use:   "rewrite <target-commit-hash>",
	Short: "Blame shifting via interactive rebase",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetHash := args[0]
		amendArgs := []string{"amend"}
		if rewriteAuthor != "" {
			amendArgs = append(amendArgs, "--author", fmt.Sprintf("'%s'", rewriteAuthor))
		}
		if rewriteVip != "" {
			amendArgs = append(amendArgs, "--vip", rewriteVip)
		}
		if cloneHash != "" {
			amendArgs = append(amendArgs, "--clone", cloneHash)
		}
		if customDate != "" {
			amendArgs = append(amendArgs, "--date", fmt.Sprintf("'%s'", customDate))
		}
		if sign {
			amendArgs = append(amendArgs, "--sign")
		}

		exePath, err := os.Executable()
		if err != nil {
			exePath = "git-forge"
		}

		fmt.Printf("Rewriting commit %s...\n", targetHash)

		execCmd := fmt.Sprintf("if [ \"$(git rev-parse HEAD)\" = \"%s\" ] || [ \"$(git rev-parse HEAD | cut -c1-7)\" = \"%s\" ]; then %s %s; fi",
			targetHash, targetHash, exePath, strings.Join(amendArgs, " "))

		rebaseArgs := []string{"rebase", "-i", targetHash + "^", "--exec", execCmd}

		os.Setenv("GIT_SEQUENCE_EDITOR", "true")

		opts := git.CommandOptions{
			DryRun:  dryRun,
			Verbose: verbose,
			NoSign:  !sign,
		}

		return git.RunGitCommand(rebaseArgs, opts)
	},
}

func init() {
	vips := strings.Join(git.GetVIPs(), ", ")
	rewriteCmd.Flags().StringVarP(&rewriteAuthor, "author", "a", "", "Explicit identity (Name <email>)")
	rewriteCmd.Flags().StringVar(&rewriteVip, "vip", "", fmt.Sprintf("Load identity from profile. Available: %s", vips))
	rootCmd.AddCommand(rewriteCmd)
}

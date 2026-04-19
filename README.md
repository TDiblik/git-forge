# git-forge

`git forge` is a CLI tool designed for educational and penetration-testing purposes to demonstrate vulnerabilities in unsigned Git repositories. It allows for the manipulation of Git identities (Author/Committer), timestamps, and GPG signatures.

## Installation

### From Source
You can use the provided `Makefile` to build and install the tool:

```bash
# Build prod and install to ~/.local/bin/git-forge so you can use it like `git forge` right away.
make use

# For dev
make build
```

## Usage
All commands support global flags like `--dry-run` (to see what would happen without altering the repo) and `--verbose` (for debugging).

### 1. Creating a Spoofed Commit
Create a new commit while impersonating a specific identity, cloning an existing one, or using a "VIP" profile. For example:

```bash
git forge commit -m "Sensitive change" --author "Linus Torvalds <torvalds@linux-foundation.org>";
git forge commit -m "Merged PR" --vip linus;
git forge commit -m "Minor fixes" --clone <commit-hash>;
git forge commit -m "Hotfix" --typo-squat "ceo@company.com";;
git forge commit -m "Team effort" --vip linus --co-author "Alice <alice@company.com>";
git forge commit -m "Backdated commit" --vip linus --date "2015-05-05 15:15:15";
```

### 2. Amending HEAD
Quickly change the identity or date of the most recent commit. For example:

```bash
git forge amend --vip linus;
git forge amend --date "2021-01-01 12:00:00";
```

### 3. Rewriting History (Blame Shifting)
Rewrite a specific commit in the past using interactive rebase logic. For example:

```bash
git forge rewrite <commit-hash> --vip linus
```

### 4. Spoofing GPG Signatures
Generate a temporary, isolated GPG key for the spoofed identity and sign the commit. This demonstrates how a "Verified" badge can be obtained if the system doesn't enforce a pre-existing trust anchor. For example:

```bash
git forge commit -m "Signed commit" --vip vitalik --sign
```

## Disclaimer
This tool is for **educational and authorized security testing only**.

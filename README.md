# git-forge

`git forge` is a CLI tool designed for educational and penetration-testing purposes to demonstrate vulnerabilities in unsigned Git repositories. It allows for the manipulation of Git identities (Author/Committer), timestamps, and GPG signatures.

## Installation

### From Source
You can use the provided `Makefile` to build and install the tool:

```bash
# Build for development
make build

# Build and install to ~/.local/bin/git-forge
make use

# Standard Go install
make install
```

## Usage

### 1. Creating a Spoofed Commit
Create a new commit while impersonating a specific identity or using a "VIP" profile.

```bash
# Using a manual author string
git-forge commit -m "Sensitive change" --author "Linus Torvalds <torvalds@linux-foundation.org>"

# Using a built-in VIP profile
git-forge commit -m "Merged PR" --vip satoshi

# Using typo-squatting to mimic an email
git-forge commit -m "Fix" --typo-squat "ceo@company.com"
```

### 2. Amending HEAD
Quickly change the identity or date of the most recent commit.

```bash
# Overwrite HEAD with a VIP identity
git-forge amend --vip linus

# Change the date of the last commit
git-forge amend --date "2021-01-01 12:00:00"
```

### 3. Rewriting History (Blame Shifting)
Rewrite a specific commit in the past using interactive rebase logic.

```bash
# Change the author of a specific commit hash
git-forge rewrite <commit-hash> --vip dhh
```

### 4. Spoofing GPG Signatures
Generate a temporary, isolated GPG key for the spoofed identity and sign the commit. This demonstrates how a "Verified" badge can be obtained if the system doesn't enforce a pre-existing trust anchor.

```bash
# Commit with a temporary GPG signature
git-forge commit -m "Signed commit" --vip vitalik --sign
```

## Disclaimer
This tool is for **educational and authorized security testing only**. Manipulating Git history and identities can be used for malicious purposes; ensure you have permission before using it on repositories you do not own.

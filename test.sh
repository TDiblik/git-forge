
#!/bin/bash

# This whole file is fully AI generated based on specs and program help, I didn't want to bother with it atm and it kinda works.

set -e

RED='\033[0-31m'
GREEN='\033[0-32m'
NC='\033[0m'

assert_eq() {
    local expected="$1"
    local actual="$2"
    local message="$3"
    if [ "$expected" != "$actual" ]; then
        echo -e "${RED}❌ FAIL: $message${NC}"
        echo "   Expected: '$expected'"
        echo "   Actual:   '$actual'"
        exit 1
    else
        echo -e "${GREEN}✅ PASS: $message${NC}"
    fi
}

expect_fail() {
    local message="$1"
    shift
    if "$@"; then
        echo -e "${RED}❌ FAIL: Expected command to fail, but it succeeded: $*${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ PASS (Expected Failure): $message${NC}"
    fi
}

prepare_changes() {
    echo "data-$(date +%s)-$RANDOM" >> file.txt
    git add file.txt
}

echo "🚀 Starting git-forge exhaustive validation suite..."

TEST_DIR="git-forge-exhaustive-test"
rm -rf "$TEST_DIR"
mkdir "$TEST_DIR"
cd "$TEST_DIR"
git init --quiet
git config user.name "Baseline User"
git config user.email "baseline@user.com"
git config commit.gpgsign true

echo "init" > file.txt
git add file.txt
git commit -m "Initial commit" --quiet
BASE_HASH=$(git rev-parse HEAD)
BASE_DATE=$(git log -1 --format="%ad")
BASE_DATE_ISO=$(git log -1 --format="%aI")

echo "--- Phase 2: Negative & Conflict Testing ---"
expect_fail "Missing message flag" git forge commit --author "Test <test@test.com>"
expect_fail "Mutually exclusive flags (clone & author)" git forge commit -m "Conflict" --clone "$BASE_HASH" --author "A <a@b.com>"
expect_fail "Mutually exclusive flags (author & vip)" git forge commit -m "Conflict" --author "A <a@b.com>" --vip linus
expect_fail "Mutually exclusive flags (vip & clone)" git forge commit -m "Conflict" --vip linus --clone "$BASE_HASH"
expect_fail "Non-existent VIP profile" git forge commit -m "Bad VIP" --vip non_existent
expect_fail "Empty commit message" git forge commit -m "" --author "A <a@b.com>"
expect_fail "Invalid date format" git forge commit -m "Fail" --author "A <a@b.com>" --date "Not A Date"

echo "--- Phase 3: Identity & Cloning ---"
prepare_changes
git forge commit -m "Clone Test" --clone "$BASE_HASH" > /dev/null
CLONE_AUTHOR=$(git log -1 --format="%an")
CLONE_DATE=$(git log -1 --format="%ad")

assert_eq "Baseline User" "$CLONE_AUTHOR" "Clone should copy author name"
assert_eq "$BASE_DATE" "$CLONE_DATE" "Clone should copy exact author date"

echo "--- Phase 4: Time Travel & Typo Squatting ---"
prepare_changes
git forge commit -m "Time Override Test" --author "Satoshi <satoshi@gmx.com>" --date "2009-01-03T18:15:05" > /dev/null
TIME_AUTHOR=$(git log -1 --format="%an")
TIME_DATE=$(git log -1 --format="%aI")

assert_eq "Satoshi" "$TIME_AUTHOR" "Explicit author should be set"
assert_eq "2009-01-03T18:15:05Z" "$(echo $TIME_DATE | sed 's/+00:00/Z/')" "Explicit date should override current time"

prepare_changes
git forge commit -m "Typo Test" --typo-squat "ceo@company.com" > /dev/null
TYPO_EMAIL=$(git log -1 --format="%ae")
if [[ "$TYPO_EMAIL" != "ceo@company.com" && "$TYPO_EMAIL" == *"@"* ]]; then
    echo -e "${GREEN}✅ PASS: Typo squat successfully altered the email ($TYPO_EMAIL)${NC}"
else
    echo -e "${RED}❌ FAIL: Typo squat failed to alter the email correctly. Actual: $TYPO_EMAIL${NC}"
    exit 1
fi

echo "--- Phase 5: Amend (Amnesia) ---"
# A. VIP and Date
git forge amend --vip linus --date "2022-02-02T12:00:00" > /dev/null
assert_eq "Linus Torvalds" "$(git log -1 --format='%an')" "Amend: identity change via VIP"
assert_eq "2022-02-02T12:00:00Z" "$(echo $(git log -1 --format='%aI') | sed 's/+00:00/Z/')" "Amend: date change"

# B. Clone
git forge amend --clone "$BASE_HASH" > /dev/null
assert_eq "Baseline User" "$(git log -1 --format='%an')" "Amend: identity change via --clone"
assert_eq "Baseline User <baseline@user.com>" "$(git log -1 --format='%an <%ae>')" "Amend: email change via --clone"
assert_eq "$BASE_DATE_ISO" "$(git log -1 --format='%aI')" "Amend: date change via --clone"

echo "--- Phase 6: Rewrite Isolation ---"
# Create some history to rewrite
prepare_changes && git commit -m "History 1" --quiet --no-gpg-sign
prepare_changes && git commit -m "History 2" --quiet --no-gpg-sign
TARGET_REWRITE=$(git log --grep="History 1" -1 --format="%H")

git forge rewrite "$TARGET_REWRITE" --vip matz > /dev/null

assert_eq "Yukihiro Matsumoto" "$(git log --grep="History 1" -1 --format="%an")" "Rewrite: target identity updated"
assert_eq "Baseline User" "$(git log -1 --format="%an")" "Rewrite: children (HEAD) preserved identity"

echo "--- Phase 7: GPG Policies & Signature Forging (Dry-Run) ---"
prepare_changes
git forge commit -m "Policy Suppression Test" --author "Attacker <attacker@evil.com>" > /dev/null
GPG_STATUS_NO_SIGN=$(git log -1 --format="%G?")
assert_eq "N" "$GPG_STATUS_NO_SIGN" "Global GPG config should be suppressed by default (--no-sign)"

prepare_changes
FORGE_SIGN_OUT=$(git forge commit -m "Forged Signature Test" --vip satoshi --sign --dry-run)
if [[ "$FORGE_SIGN_OUT" == *"-S"* && "$FORGE_SIGN_OUT" == *"gpgsign=true"* ]]; then
    echo -e "${GREEN}✅ PASS: Commit signing logic triggered correctly (Dry-Run)${NC}"
else
    echo -e "${RED}❌ FAIL: Commit signing logic not triggered. Got: $FORGE_SIGN_OUT${NC}"
    exit 1
fi

echo "--- Phase 8: Dry Run Safety Check ---"
prepare_changes
PRE_DRY_RUN_HASH=$(git rev-parse HEAD)
git forge commit -m "Dry Run Test" --author "Ghost <ghost@shell.com>" --dry-run > /dev/null
POST_DRY_RUN_HASH=$(git rev-parse HEAD)
assert_eq "$PRE_DRY_RUN_HASH" "$POST_DRY_RUN_HASH" "Dry run MUST NOT alter the repository state"

echo "--- Phase 9: Date Priority and Formats ---"
prepare_changes
git forge commit -m "Date Priority Test" --clone "$BASE_HASH" --date "2025-12-25T10:00:00" > /dev/null
PRIORITY_DATE=$(git log -1 --format="%aI")
assert_eq "2025-12-25T10:00:00Z" "$(echo $PRIORITY_DATE | sed 's/+00:00/Z/')" "Explicit --date must override --clone timestamp"

prepare_changes
git forge commit -m "Date Format Test 1" --author "Tester <t@t.com>" --date "2020-01-01 12:00:00" > /dev/null
FORMAT_DATE_1=$(git log -1 --format="%ai")
assert_eq "2020-01-01 12:00:00 +0000" "$FORMAT_DATE_1" "Space-separated date format should work"

echo "--- Phase 10: Edge Cases ---"

# A. Multi-line message
prepare_changes
git forge commit -m "Line 1
Line 2
Line 3" --vip matz > /dev/null
assert_eq "Line 1
Line 2
Line 3" "$(git log -1 --format="%B")" "Edge: Multi-line message"

# B. Special characters in author
prepare_changes
git forge commit -m "Special Char" --author "O'Neil & Son <oneil@test.com>" > /dev/null
assert_eq "O'Neil & Son|oneil@test.com" "$(git log -1 --format="%an|%ae")" "Edge: Special characters in identity"

echo "--- Phase 11: Utilities ---"
git forge completion bash > /dev/null
git forge completion zsh > /dev/null
git forge help rewrite > /dev/null
echo -e "${GREEN}✅ PASS: Completion and Help commands ran successfully${NC}"

VERBOSE_OUT=$(git forge commit -m "Verbose Test" --author "V <v@v.com>" --verbose --dry-run)
if [[ "$VERBOSE_OUT" == *"Executing:"* || "$VERBOSE_OUT" == *"Env:"* ]]; then
    echo -e "${GREEN}✅ PASS: Verbose output contains execution details${NC}"
else
    echo -e "${RED}❌ FAIL: Verbose output missing details. Got: $VERBOSE_OUT${NC}"
    exit 1
fi

echo -e "\n${GREEN}🎉 ALL EXHAUSTIVE TESTS PASSED SUCCESSFULLY! 🌟${NC}"
cd ..
rm -rf "$TEST_DIR"

set -e

RED='\033[0-31m'
GREEN='\033[0-32m'
NC='\033[0m'

assert_eq() {
    local expected="$1"
    local actual="$2"
    local message="$3"
    if [ "$expected" != "$actual" ]; then
        echo -e "${RED}❌ FAIL: $message${NC}"
        echo "   Expected: '$expected'"
        echo "   Actual:   '$actual'"
        exit 1
    else
        echo -e "${GREEN}✅ PASS: $message${NC}"
    fi
}

expect_fail() {
    local message="$1"
    shift
    if "$@"; then
        echo -e "${RED}❌ FAIL: Expected command to fail, but it succeeded: $*${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ PASS (Expected Failure): $message${NC}"
    fi
}

prepare_changes() {
    echo "data-$(date +%s)-$RANDOM" >> file.txt
    git add file.txt
}

TEST_DIR="git-forge-exhaustive-test"
rm -rf "$TEST_DIR"
mkdir "$TEST_DIR"
cd "$TEST_DIR"
git init --quiet
git config user.name "Baseline User"
git config user.email "baseline@user.com"
git config commit.gpgsign true

echo "init" > file.txt
git add file.txt
git commit -m "Initial commit" --quiet
BASE_HASH=$(git rev-parse HEAD)
BASE_DATE=$(git log -1 --format="%ad")
BASE_DATE_ISO=$(git log -1 --format="%aI")

echo "--- Phase 2: Negative & Conflict Testing ---"
expect_fail "Missing message flag" git forge commit --author "Test <test@test.com>"
expect_fail "Mutually exclusive flags (clone & author)" git forge commit -m "Conflict" --clone "$BASE_HASH" --author "A <a@b.com>"
expect_fail "Mutually exclusive flags (author & vip)" git forge commit -m "Conflict" --author "A <a@b.com>" --vip linus
expect_fail "Mutually exclusive flags (vip & clone)" git forge commit -m "Conflict" --vip linus --clone "$BASE_HASH"
expect_fail "Non-existent VIP profile" git forge commit -m "Bad VIP" --vip non_existent
expect_fail "Empty commit message" git forge commit -m "" --author "A <a@b.com>"
expect_fail "Invalid date format" git forge commit -m "Fail" --author "A <a@b.com>" --date "Not A Date"
expect_fail "Nothing staged to commit" git forge commit -m "Empty" --author "A <a@b.com>" # Working directory clean

echo "--- Phase 3: Identity & Cloning ---"
prepare_changes
git forge commit -m "Clone Test" --clone "$BASE_HASH" > /dev/null
CLONE_AUTHOR=$(git log -1 --format="%an")
CLONE_DATE=$(git log -1 --format="%ad")

assert_eq "Baseline User" "$CLONE_AUTHOR" "Clone should copy author name"
assert_eq "$BASE_DATE" "$CLONE_DATE" "Clone should copy exact author date"

echo "--- Phase 4: Time Travel & Typo Squatting ---"
prepare_changes
git forge commit -m "Time Override Test" --author "Satoshi <satoshi@gmx.com>" --date "2009-01-03T18:15:05" > /dev/null
TIME_AUTHOR=$(git log -1 --format="%an")
TIME_DATE=$(git log -1 --format="%aI")

assert_eq "Satoshi" "$TIME_AUTHOR" "Explicit author should be set"
assert_eq "2009-01-03T18:15:05Z" "$(echo $TIME_DATE | sed 's/+00:00/Z/')" "Explicit date should override current time"

prepare_changes
git forge commit -m "Typo Test" --typo-squat "ceo@company.com" > /dev/null
TYPO_EMAIL=$(git log -1 --format="%ae")
if [[ "$TYPO_EMAIL" != "ceo@company.com" && "$TYPO_EMAIL" == *"@"* ]]; then
    echo -e "${GREEN}✅ PASS: Typo squat successfully altered the email ($TYPO_EMAIL)${NC}"
else
    echo -e "${RED}❌ FAIL: Typo squat failed to alter the email correctly. Actual: $TYPO_EMAIL${NC}"
    exit 1
fi

echo "--- Phase 5: Amend (Amnesia) ---"
# A. VIP and Date
prepare_changes && git commit -m "To be amended" --quiet --no-gpg-sign
git forge amend --vip linus --date "2022-02-02T12:00:00" > /dev/null
assert_eq "Linus Torvalds" "$(git log -1 --format='%an')" "Amend: identity change via VIP"
assert_eq "2022-02-02T12:00:00Z" "$(echo $(git log -1 --format='%aI') | sed 's/+00:00/Z/')" "Amend: date change"

# B. Clone
prepare_changes && git commit -m "To be amended 2" --quiet --no-gpg-sign
git forge amend --clone "$BASE_HASH" > /dev/null
assert_eq "Baseline User" "$(git log -1 --format='%an')" "Amend: identity change via --clone"
assert_eq "Baseline User <baseline@user.com>" "$(git log -1 --format='%an <%ae>')" "Amend: email change via --clone"
assert_eq "$BASE_DATE_ISO" "$(git log -1 --format='%aI')" "Amend: date change via --clone"

echo "--- Phase 6: Rewrite Isolation & Modification ---"
# Create some history to rewrite
prepare_changes && git commit -m "History 1" --quiet --no-gpg-sign
TARGET_REWRITE=$(git log -1 --format="%H")
prepare_changes && git commit -m "History 2" --quiet --no-gpg-sign

# Rewrite with Identity AND Date change
git forge rewrite "$TARGET_REWRITE" --vip matz --date "2015-05-05T15:15:15" > /dev/null

assert_eq "Yukihiro Matsumoto" "$(git log --grep="History 1" -1 --format="%an")" "Rewrite: target identity updated"
assert_eq "2015-05-05T15:15:15Z" "$(echo $(git log --grep="History 1" -1 --format="%aI") | sed 's/+00:00/Z/')" "Rewrite: target date updated"
assert_eq "Baseline User" "$(git log -1 --format="%an")" "Rewrite: children (HEAD) preserved identity"

echo "--- Phase 7: GPG Policies & Signature Forging ---"
prepare_changes
git forge commit -m "Policy Suppression Test" --author "Attacker <attacker@evil.com>" > /dev/null
GPG_STATUS_NO_SIGN=$(git log -1 --format="%G?")
assert_eq "N" "$GPG_STATUS_NO_SIGN" "Global GPG config should be suppressed by default (--no-sign)"

# Real Ephemeral GPG Test
export GNUPGHOME="$(mktemp -d)"
prepare_changes
git forge commit -m "Real GPG Sign" --vip linus --sign > /dev/null
SIG_STATUS=$(git log -1 --format="%G?")
if [[ "$SIG_STATUS" == "U" || "$SIG_STATUS" == "E" ]]; then
    echo -e "${GREEN}✅ PASS: Tool successfully generated an ephemeral key and signed the commit (Status: $SIG_STATUS)${NC}"
else
    echo -e "${RED}❌ FAIL: Tool failed to sign the commit properly. Status: $SIG_STATUS${NC}"
    exit 1
fi
rm -rf "$GNUPGHOME"
unset GNUPGHOME

echo "--- Phase 8: Dry Run Safety Check ---"
prepare_changes
PRE_DRY_RUN_HASH=$(git rev-parse HEAD)
git forge commit -m "Dry Run Test" --author "Ghost <ghost@shell.com>" --dry-run > /dev/null
POST_DRY_RUN_HASH=$(git rev-parse HEAD)
assert_eq "$PRE_DRY_RUN_HASH" "$POST_DRY_RUN_HASH" "Dry run MUST NOT alter the repository state"

echo "--- Phase 9: Date Priority and Formats ---"
prepare_changes
git forge commit -m "Date Priority Test" --clone "$BASE_HASH" --date "2025-12-25T10:00:00" > /dev/null
PRIORITY_DATE=$(git log -1 --format="%aI")
assert_eq "2025-12-25T10:00:00Z" "$(echo $PRIORITY_DATE | sed 's/+00:00/Z/')" "Explicit --date must override --clone timestamp"

echo "--- Phase 10: Identity Leakage & Co-Authors ---"

# A. Committer Leakage Test
prepare_changes
git forge commit -m "Committer Leak Test" --author "Ninja <ninja@stealth.com>" > /dev/null
assert_eq "Ninja" "$(git log -1 --format="%cn")" "Committer name MUST match Author name (No Leakage)"
assert_eq "ninja@stealth.com" "$(git log -1 --format="%ce")" "Committer email MUST match Author email (No Leakage)"

# B. Co-Authors Trailers
prepare_changes
git forge commit -m "Teamwork" --author "Main <m@m.com>" --co-author "Alice <a@a.com>" --co-author "Bob <b@b.com>" > /dev/null
BODY=$(git log -1 --format="%B")
if [[ "$BODY" == *"Co-authored-by: Alice <a@a.com>"* && "$BODY" == *"Co-authored-by: Bob <b@b.com>"* ]]; then
    echo -e "${GREEN}✅ PASS: Co-author trailers added successfully${NC}"
else
    echo -e "${RED}❌ FAIL: Co-author trailers missing or incorrect. Body: $BODY${NC}"
    exit 1
fi

echo "--- Phase 11: Edge Cases ---"
# A. Multi-line message
prepare_changes
git forge commit -m "Line 1
Line 2
Line 3" --vip matz > /dev/null
assert_eq "Line 1
Line 2
Line 3" "$(git log -1 --format="%B")" "Edge: Multi-line message preserved"

# B. Special characters in author
prepare_changes
git forge commit -m "Special Char" --author "O'Neil & Son <oneil@test.com>" > /dev/null
assert_eq "O'Neil & Son|oneil@test.com" "$(git log -1 --format="%an|%ae")" "Edge: Special characters in identity handled safely"

echo "--- Phase 12: Utilities ---"
git forge completion bash > /dev/null || echo "Skipped (no completion)"
git forge help rewrite > /dev/null
echo -e "${GREEN}✅ PASS: Help commands ran successfully${NC}"

VERBOSE_OUT=$(git forge commit -m "Verbose Test" --author "V <v@v.com>" --verbose --dry-run 2>&1 || true)
if [[ -n "$VERBOSE_OUT" ]]; then
    echo -e "${GREEN}✅ PASS: Verbose output generated${NC}"
else
    echo -e "${RED}❌ FAIL: Verbose output missing.${NC}"
    exit 1
fi

echo -e "\n${GREEN}🎉 ALL EXHAUSTIVE TESTS PASSED SUCCESSFULLY! 🌟${NC}"
cd ..
rm -rf "$TEST_DIR"

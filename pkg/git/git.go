package git

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Identity struct {
	Name  string
	Email string
	Date  string
}

type CommandOptions struct {
	Identity   *Identity
	DryRun     bool
	Verbose    bool
	SigningKey string
	NoSign     bool
}

func ResolveFromHash(hash string) (*Identity, error) {
	cmd := exec.Command("git", "show", "-s", "--format=%an|%ae|%ai", hash)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hash %s: %v: %s", hash, err, string(out))
	}

	parts := strings.Split(strings.TrimSpace(string(out)), "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected output from git show: %s", string(out))
	}

	return &Identity{
		Name:  parts[0],
		Email: parts[1],
		Date:  parts[2],
	}, nil
}

func ParseAuthor(authorStr string) (*Identity, error) {
	if !strings.Contains(authorStr, "<") || !strings.Contains(authorStr, ">") {
		return nil, fmt.Errorf("invalid author format, use 'Name <email>'")
	}
	parts := strings.Split(authorStr, "<")
	name := strings.TrimSpace(parts[0])
	email := strings.Trim(parts[1], "> ")
	return &Identity{Name: name, Email: email}, nil
}

func RunGitCommand(args []string, opts CommandOptions) error {
	var finalArgs []string

	if opts.SigningKey != "" {
		finalArgs = append([]string{"-c", fmt.Sprintf("user.signingkey=%s", opts.SigningKey)}, args...)
	} else if opts.NoSign {
		isCommitCmd := false
		for _, arg := range args {
			if arg == "commit" || arg == "rebase" || arg == "amend" {
				isCommitCmd = true
				break
			}
		}
		if isCommitCmd {
			finalArgs = append(args, "--no-sign")
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}

	cmd := exec.Command("git", finalArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if opts.Identity != nil {
		env := os.Environ()
		env = append(env, fmt.Sprintf("GIT_AUTHOR_NAME=%s", opts.Identity.Name))
		env = append(env, fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", opts.Identity.Email))
		env = append(env, fmt.Sprintf("GIT_COMMITTER_NAME=%s", opts.Identity.Name))
		env = append(env, fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", opts.Identity.Email))

		if opts.Identity.Date != "" {
			env = append(env, fmt.Sprintf("GIT_AUTHOR_DATE=%s", opts.Identity.Date))
			env = append(env, fmt.Sprintf("GIT_COMMITTER_DATE=%s", opts.Identity.Date))
		}
		cmd.Env = env
	}

	if opts.DryRun {
		fmt.Printf("[DRY-RUN] Executing: git %s\n", strings.Join(finalArgs, " "))
		if opts.Identity != nil {
			fmt.Printf("[DRY-RUN] Envs: AUTHOR_NAME=%s, AUTHOR_EMAIL=%s, DATE=%s\n", opts.Identity.Name, opts.Identity.Email, opts.Identity.Date)
		}
		return nil
	}

	if opts.Verbose {
		fmt.Printf("Executing: git %s\n", strings.Join(finalArgs, " "))
	}

	return cmd.Run()
}

func TypoSquat(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}
	parts := strings.Split(email, "@")
	if len(parts[0]) < 2 {
		return email
	}
	local := parts[0]
	squatted := local[:len(local)-2] + string(local[len(local)-1]) + string(local[len(local)-2])
	return squatted + "@" + parts[1]
}

func ResolveVIP(profile string) (*Identity, error) {
	id, ok := vip_profiles[strings.ToLower(profile)]
	if !ok {
		return nil, fmt.Errorf("unknown VIP profile: %s", profile)
	}
	return id, nil
}

func GetVIPs() []string {
	keys := make([]string, 0, len(vip_profiles))
	for k := range vip_profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var vip_profiles = map[string]*Identity{
	"linus":               {Name: "Linus Torvalds", Email: "torvalds@linux-foundation.org"},
	"satoshi":             {Name: "Satoshi Nakamoto", Email: "satoshin@gmx.com"},
	"guido":               {Name: "Guido van Rossum", Email: "guido@python.org"},
	"dhh":                 {Name: "David Heinemeier Hansson", Email: "david@loudpost.com"},
	"antirez":             {Name: "Salvatore Sanfilippo", Email: "antirez@gmail.com"},
	"robpike":             {Name: "Rob Pike", Email: "r@google.com"},
	"ken":                 {Name: "Ken Thompson", Email: "ken@google.com"},
	"matz":                {Name: "Yukihiro Matsumoto", Email: "matz@ruby-lang.org"},
	"vladimir":            {Name: "Wladimir J. van der Laan", Email: "laanwj@protonmail.com"},
	"rms":                 {Name: "Richard Stallman", Email: "rms@gnu.org"},
	"gkh":                 {Name: "Greg Kroah-Hartman", Email: "gregkh@linuxfoundation.org"},
	"stroustrup":          {Name: "Bjarne Stroustrup", Email: "bjarne@stroustrup.com"},
	"gosling":             {Name: "James Gosling", Email: "james.gosling@sun.com"},
	"brendan":             {Name: "Brendan Eich", Email: "brendan@mozilla.org"},
	"rasmus":              {Name: "Rasmus Lerdorf", Email: "rasmus@php.net"},
	"taylor":              {Name: "Taylor Otwell", Email: "taylor@laravel.com"},
	"evanyou":             {Name: "Evan You", Email: "evan@vuejs.org"},
	"gaearon":             {Name: "Dan Abramov", Email: "dan.abramov@gmail.com"},
	"mitchellh":           {Name: "Mitchell Hashimoto", Email: "m@mitchellh.com"},
	"shykes":              {Name: "Solomon Hykes", Email: "solomon@docker.com"},
	"kelsey":              {Name: "Kelsey Hightower", Email: "kelsey.hightower@gmail.com"},
	"clattner":            {Name: "Chris Lattner", Email: "clattner@nondot.org"},
	"graydon":             {Name: "Graydon Hoare", Email: "graydon@mozilla.com"},
	"bradfitz":            {Name: "Brad Fitzpatrick", Email: "brad@danga.com"},
	"vitalik":             {Name: "Vitalik Buterin", Email: "vitalik@ethereum.org"},
	"carmack":             {Name: "John Carmack", Email: "john.carmack@oculus.com"},
	"gates":               {Name: "Bill Gates", Email: "billg@microsoft.com"},
	"zuck":                {Name: "Mark Zuckerberg", Email: "zuck@fb.com"},
	"musk":                {Name: "Elon Musk", Email: "elon@spacex.com"},
	"woz":                 {Name: "Steve Wozniak", Email: "woz@apple.com"},
	"bezos":               {Name: "Jeff Bezos", Email: "jeff@amazon.com"},
	"cook":                {Name: "Tim Cook", Email: "tcook@apple.com"},
	"pichai":              {Name: "Sundar Pichai", Email: "sundar@google.com"},
	"nadella":             {Name: "Satya Nadella", Email: "satyan@microsoft.com"},
	"swartz":              {Name: "Aaron Swartz", Email: "aaron@aaronsw.com"},
	"snowden":             {Name: "Edward Snowden", Email: "snowden@protonmail.com"},
	"assange":             {Name: "Julian Assange", Email: "proff@iq.org"},
	"ry":                  {Name: "Ryan Dahl", Email: "ry@tinyclouds.org"},
	"tj":                  {Name: "TJ Holowaychuk", Email: "tj@vision-media.ca"},
	"rsc":                 {Name: "Russ Cox", Email: "rsc@golang.org"},
	"fabrice":             {Name: "Fabrice Bellard", Email: "fabrice@bellard.org"},
	"jeresig":             {Name: "John Resig", Email: "jeresig@gmail.com"},
	"sindre":              {Name: "Sindre Sorhus", Email: "sindresorhus@gmail.com"},
	"tenderlove":          {Name: "Aaron Patterson", Email: "aaron@tenderlovemaking.com"},
	"mitsuhiko":           {Name: "Armin Ronacher", Email: "armin.ronacher@active-4.com"},
	"jose":                {Name: "José Valim", Email: "jose.valim@dashbit.co"},
	"bram":                {Name: "Bram Moolenaar", Email: "Bram@moolenaar.net"},
	"crockford":           {Name: "Douglas Crockford", Email: "douglas@crockford.com"},
	"jordwalke":           {Name: "Jordan Walke", Email: "jordwalke@gmail.com"},
	"richharris":          {Name: "Rich Harris", Email: "richard.a.harris@gmail.com"},
	"paulirish":           {Name: "Paul Irish", Email: "paul.irish@gmail.com"},
	"addy":                {Name: "Addy Osmani", Email: "addyosmani@gmail.com"},
	"hixie":               {Name: "Ian Hickson", Email: "ian@hixie.ch"},
	"wycats":              {Name: "Yehuda Katz", Email: "wycats@gmail.com"},
	"sophie":              {Name: "Sophie Alpert", Email: "sophie@sophiebits.com"},
	"kentcdodds":          {Name: "Kent C. Dodds", Email: "kent@doddsfamily.us"},
	"rauchg":              {Name: "Guillermo Rauch", Email: "rauchg@gmail.com"},
	"fireship":            {Name: "Jeff Delaney", Email: "jeff@fireship.io"},
	"primeagen":           {Name: "Michael Paulson", Email: "the.primeagen@gmail.com"},
	"theo":                {Name: "Theo Browne", Email: "ping@t3.gg"},
	"mkbhd":               {Name: "Marques Brownlee", Email: "marques@mkbhd.com"},
	"linustech":           {Name: "Linus Sebastian", Email: "info@linusmediagroup.com"},
	"tomscott":            {Name: "Tom Scott", Email: "tom@tomscott.com"},
	"michaelreeves":       {Name: "Michael Reeves", Email: "michaelreeves@wmeagency.com"},
	"programmingwithmosh": {Name: "Mosh Hamedani", Email: "programmingwithmosh@gmail.com"},
	"ceo":                 {Name: "The CEO", Email: "ceo@company.com"},
	"cto":                 {Name: "The CTO", Email: "cto@company.com"},
}

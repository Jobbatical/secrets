// `options` handles reading in user input. Does not have logic for detecting smart defaults.
package options

import (
	"errors"
	"flag"
	"fmt"
	"jobbatical/secrets/pkg/utils"
	"os"
	"strconv"
	"path/filepath"
	"strings"
)

const Usage string = "Usage secrets <open|seal> [<file path>...] [--dry-run] [--verbose] [--root <project root>] [--key <encryption key name>] [--open-all]"
const EncryptCmd string = "seal"
const DecryptCmd string = "open"

var ExpectedOrganization string
var ExpectedRepoHost string
var KeyRing string
var KeyLocation string

var Concurrency int
var DryRun bool
var Key string
var OpenAll bool
var ProjectRoot string
var Verbose bool
var Cmd string
var Files []string

func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func popCommand(args []string) (string, []string, error) {
	for i, a := range args {
		if i == 0 {
			continue
		}
		if !strings.HasPrefix(a, "-") {
			return a, Remove(args, i), nil
		} else {
			continue
		}
	}
	return "", args, errors.New("command not found")
}

func popFiles(args []string) ([]string, []string, error) {
	var (
		file string
		err  error
	)
	files := make([]string, 0, 1)

	for {
		file, os.Args, err = popCommand(os.Args)
		if err != nil {
			break
		}
		absolutePath, err := filepath.Abs(file)
		if err != nil {
			return files, os.Args, err
		}
		files = append(files, absolutePath)
	}

	return files, os.Args, nil
}

func readRequiredConfig(v *string, argFlag string, envName string, desc string) func() {
	flag.StringVar(v, argFlag, os.Getenv(envName), desc)
	return func() {
		if len(*v) == 0 {
			fmt.Fprintf(os.Stderr, "Missing required configuration: %s\nPlease set %s environment var or pass --%s flag\n", desc, envName, argFlag)
			os.Exit(2)
		}
	}
}

func getConcurrency() (int) {
	concurrencyEnv := os.Getenv("SECRETS_CONCURRENCY")
	if (len(concurrencyEnv) == 0) {
		return 5
	}
	i, err := strconv.Atoi(concurrencyEnv)
	utils.ExitIfError(err)
	return i
}

func init() {
	var err error

	Cmd, os.Args, err = popCommand(os.Args)
	if err != nil {
		utils.ErrPrintln("Error: %s\n%s", err, Usage)
		os.Exit(1)
	}

	Files, os.Args, err = popFiles(os.Args)
	utils.ExitIfError(err)

	// Execution options
	flag.BoolVar(&DryRun, "dry-run", false, "Skip calls to GCP")
	flag.StringVar(&Key, "key", "", "Key to use")
	flag.BoolVar(&OpenAll, "open-all", false, "Opens all .enc files within the repository")
	flag.StringVar(&ProjectRoot, "root", "", "Project root folder(name will be used as key name)")
	flag.BoolVar(&Verbose, "verbose", false, "Log debug info")
	flag.IntVar(&Concurrency, "j", getConcurrency(), "Number of concurrent calls to the service provider")

	// Configuration
	defer readRequiredConfig(&ExpectedOrganization, "org", "SECRETS_ORG", "Expected organization of the repo")()
	defer readRequiredConfig(&ExpectedRepoHost, "repo-host", "SECRETS_REPO_HOST", "Expected host for the repo")()
	defer readRequiredConfig(&KeyRing, "key-ring", "SECRETS_KEY_RING", "The key ring to use for encryption")()
	defer readRequiredConfig(&KeyLocation, "key-location", "SECRETS_KEY_LOCATION", "The location of the key ring")()

	flag.Parse()
}

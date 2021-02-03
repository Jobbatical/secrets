package main

import (
	"errors"
	"fmt"
	"jobbatical/secrets/pkg/git"
	"jobbatical/secrets/pkg/kms"
	"jobbatical/secrets/pkg/log"
	"jobbatical/secrets/pkg/options"
	"jobbatical/secrets/pkg/utils"
	"os"
	"path/filepath"
)

var verbose bool = options.Verbose
var projectRoot string = options.ProjectRoot
var key string = options.Key

func isProjectRoot(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		return false
	}
	return info.IsDir()
}

func findProjectRoot(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	nextPath := filepath.Join(path, "..")
	if path == nextPath {
		return path, errors.New("not in project. Run the script inside a project folder(git repo) or provide it as an argument")
	}
	if isProjectRoot(path) {
		return path, nil
	}
	return findProjectRoot(nextPath)
}

func getKeyName(projectRoot string) string {
	repo, err := git.GetProjectRepo(projectRoot, options.ExpectedRepoHost, options.ExpectedOrganization)
	if err == nil {
		return repo
	}
	return filepath.Base(projectRoot)
}

func main() {
	log.PrintDebugln("%s", os.Args)

	if projectRoot == "" {
		projectRoot, _ = findProjectRoot(".")
	}

	if key == "" {
		key = getKeyName(projectRoot)
	}

	log.PrintDebugln("dry run: %t", options.DryRun)
	log.PrintDebugln("concurrency: %t", options.Concurrency)
	log.PrintDebugln("options.ExpectedOrganization: %s", options.ExpectedOrganization)
	log.PrintDebugln("options.ExpectedRepoHost: %s", options.ExpectedRepoHost)
	log.PrintDebugln("keyRing: %s", options.KeyRing)
	log.PrintDebugln("key: %s", key)
	log.PrintDebugln("location: %s", options.KeyLocation)
	log.PrintDebugln("project root: %s", projectRoot)
	log.PrintDebugln("cmd: %s", options.Cmd)
	log.PrintDebugln("files: %s (%d)", options.Files, len(options.Files))

	if options.Cmd == options.EncryptCmd {
		if len(options.Files) == 0 {
			options.Files, _ = utils.FindUnencryptedFiles(projectRoot)
		}

		utils.Concurrently(options.Concurrency, options.Files, func(path string) {
			fmt.Printf("encrypting %s\n", path)
			utils.ExitIfError(kms.Encrypt(key, path))
			err := git.AddToIgnored(projectRoot, path)
			if err == git.ErrFileAlreadyTracked {
				utils.ErrPrintln("Warning: plain-text file already checked in: %s", path)
				return
			}
			utils.ExitIfError(err)
		})

		os.Exit(0)
	}
	if options.Cmd == options.DecryptCmd {
		if len(options.Files) == 0 {
			options.Files, _ = utils.FindEncryptedFiles(options.OpenAll, projectRoot)
		}

		utils.Concurrently(options.Concurrency, options.Files, func(path string) {
			fmt.Printf("decrypting %s\n", path)
			err := kms.Decrypt(key, path)
			utils.ExitIfError(err)
		})

		os.Exit(0)
	}
	utils.ErrPrintln("Unknown command: %s\n%s", options.Cmd, options.Usage)
	os.Exit(1)
}

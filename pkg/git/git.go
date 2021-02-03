package git

import (
	"errors"
	"jobbatical/secrets/pkg/log"
	"jobbatical/secrets/pkg/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"regexp"
	"fmt"
)

var ErrFileAlreadyTracked = errors.New("file already tracked")

func isTracked(projectRoot string, filePath string) (bool, error) {
	_, _, _, err := utils.RunCommand(
		"git",
		"-C", projectRoot,
		"ls-files", "--error-unmatch", filePath,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func isIgnored(projectRoot string, filePath string) (bool, error) {
	_, stdOut, _, err := utils.RunCommand(
		"git",
		"-C", projectRoot,
		"check-ignore", filePath,
	)
	if err != nil {
		return false, err
	}
	return (strings.TrimSpace(stdOut) == filePath), nil
}

func appendToFile(filePath string, line string) error {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil
}

func AddToIgnored(projectRoot string, fileToIgnore string) error {
	relativePath, err := filepath.Rel(projectRoot, fileToIgnore)
	if err != nil {
		return err
	}

	isTracked, err := isTracked(projectRoot, relativePath)
	if isTracked {
		log.PrintDebugln("NOT appending %s to gitignore because it's already tracked", fileToIgnore)
		return ErrFileAlreadyTracked
	}
	isIgnored, err := isIgnored(projectRoot, fileToIgnore)
	if isIgnored {
		log.PrintDebugln("NOT appending %s to gitignore because it's already ignored", fileToIgnore)
		return nil
	}
	return appendToFile(path.Join(projectRoot, ".gitignore"), relativePath)
}

func GetProjectRepo(projectRoot string, repoHost string, org string) (string, error) {
	_, stdOut, _, err := utils.RunCommand("git", "-C", projectRoot, "remote", "-v")
	if err != nil {
		return "", err
	}
	example := fmt.Sprintf("git@%s:%s/<project name>.git", repoHost, org)
	re := regexp.MustCompile("(?i)" + repoHost + `:([^/]*)/([^/\.]*)\.git`)
	matches := re.FindStringSubmatch(stdOut)
	if len(matches) == 3 {
		org := matches[1]
		project := matches[2]

		if strings.ToLower(org) == org {
			return project, nil
		}

		return "", fmt.Errorf(
			`%s not a %s project in %s: expecting a remote %s, got %s in %s`,
			projectRoot,
			org,
			repoHost,
			example,
			project,
			org,
		)
	}
	return "", fmt.Errorf(
		`%s not a project in %s: expecting a remote %s`,
		projectRoot,
		repoHost,
		example,
	)
}

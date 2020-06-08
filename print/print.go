package print

import (
	"fmt"
	"git-get/git"
	"path/filepath"
	"strings"
)

type Printer interface {
	Print(root string, repos []*git.Repo) string
}

type FlatPrinter struct{}

func NewFlatPrinter() *FlatPrinter {
	return &FlatPrinter{}
}

func (p *FlatPrinter) Print(root string, repos []*git.Repo) string {
	val := root

	for _, repo := range repos {
		path := strings.TrimPrefix(repo.Path, root)
		path = strings.Trim(path, string(filepath.Separator))

		val += fmt.Sprintf("\n%s %s", path, renderWorktreeStatus(repo))

		for _, branch := range repo.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == repo.Status.CurrentBranch {
				continue
			}

			indent := strings.Repeat(" ", len(path))
			val += fmt.Sprintf("\n%s %s", indent, renderBranchStatus(branch))
		}
	}

	return val
}

const (
	ColorRed    = "\033[1;31m%s\033[0m"
	ColorGreen  = "\033[1;32m%s\033[0m"
	ColorBlue   = "\033[1;34m%s\033[0m"
	ColorYellow = "\033[1;33m%s\033[0m"
)

func renderWorktreeStatus(repo *git.Repo) string {
	clean := true
	var status []string

	// if current branch status can't be found it's probably a detached head
	// TODO: what if current HEAD points to a tag?
	if current := repo.CurrentBranchStatus(); current == nil {
		status = append(status, fmt.Sprintf(ColorYellow, repo.Status.CurrentBranch))
	} else {
		status = append(status, renderBranchStatus(current))
	}

	// TODO: this is ugly
	// unset clean flag to use it to render braces around worktree status and remove "ok" from branch status if it's there
	if repo.Status.HasUncommittedChanges || repo.Status.HasUntrackedFiles {
		clean = false
	}

	if !clean {
		status[len(status)-1] = strings.TrimSuffix(status[len(status)-1], git.StatusOk)
		status = append(status, "[")
	}

	if repo.Status.HasUntrackedFiles {
		status = append(status, fmt.Sprintf(ColorRed, git.StatusUntracked))
	}

	if repo.Status.HasUncommittedChanges {
		status = append(status, fmt.Sprintf(ColorRed, git.StatusUncommitted))
	}

	if !clean {
		status = append(status, "]")
	}

	return strings.Join(status, " ")
}

func renderBranchStatus(branch *git.BranchStatus) string {
	// ok indicates that the branch has upstream and is not ahead or behind it
	ok := true
	var status []string

	status = append(status, fmt.Sprintf(ColorBlue, branch.Name))

	if branch.Upstream == "" {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, git.StatusNoUpstream))
	}

	if branch.NeedsPull {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, git.StatusBehind))
	}

	if branch.NeedsPush {
		ok = false
		status = append(status, fmt.Sprintf(ColorYellow, git.StatusAhead))
	}

	if ok {
		status = append(status, fmt.Sprintf(ColorGreen, git.StatusOk))
	}

	return strings.Join(status, " ")
}

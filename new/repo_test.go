package new

import "testing"

func TestRepoEmpty(t *testing.T) {
	repo := newTestRepo(t)

	wt, err := repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	if !status.IsClean() {
		t.Errorf("Empty repo should be clean")
	}
}

func TestRepoWithUntrackedFile(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "file")

	wt, err := repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	if status.IsClean() {
		t.Errorf("Repo with untracked file should not be clean")
	}

	if !status.IsUntracked("file") {
		t.Errorf("New file should be untracked")
	}
}

func TestRepoWithStagedFile(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "file")
	stageFile(t, repo, "file")

	wt, err := repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	if status.IsClean() {
		t.Errorf("Repo with staged file should not be clean")
	}

	if status.IsUntracked("file") {
		t.Errorf("Staged file should not be untracked")
	}
}

func TestRepoWithSingleCommit(t *testing.T) {
	repo := newTestRepo(t)
	createFile(t, repo, "file")
	stageFile(t, repo, "file")
	createCommit(t, repo, "Initial commit")

	wt, err := repo.Worktree()
	checkFatal(t, err)

	status, err := wt.Status()
	if !status.IsClean() {
		t.Errorf("Repo with committed file should be clean")
	}

	if status.IsUntracked("file") {
		t.Errorf("Committed file should not be untracked")
	}
}
func TestStatusWithModifiedFile(t *testing.T) {
	//todo modified but not staged
}

func TestStatusWithUntrackedButIgnoredFile(t *testing.T) {
	//todo
}
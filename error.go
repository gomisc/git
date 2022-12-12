package git

import (
	"git.eth4.dev/golibs/errors"
)

const (
	ErrWrongRepoPath    = errors.Const("wrong repo path")
	ErrRepoPathNotEmpty = errors.Const("repo path directory is not empty")
)

package git

import (
	"gopkg.in/gomisc/errors.v1"
)

const (
	ErrWrongRepoPath    = errors.Const("wrong repo path")
	ErrRepoPathNotEmpty = errors.Const("repo path directory is not empty")
)

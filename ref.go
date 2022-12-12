package git

import (
	"github.com/go-git/go-git/v5/plumbing"
)

type gitReference struct {
	ref *plumbing.Reference
}

func (r *gitReference) Name() string {
	return r.ref.Name().Short()
}

func (r *gitReference) Hash() string {
	return r.ref.Hash().String()
}

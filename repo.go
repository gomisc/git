package git

import (
	"sync"
	"time"

	"git.eth4.dev/golibs/errors"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type gitRepository struct {
	path   string
	config *config.Config
	auth   transport.AuthMethod

	sync.RWMutex
	repo *gogit.Repository
}

// Pull пуллит репу
func (r *gitRepository) Pull(options ...Option[gogit.PullOptions]) error {
	r.Lock()
	defer r.Unlock()

	tree, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "getting work tree")
	}

	opts := processOptions(options...)
	opts.Auth = r.auth

	if opts.RemoteName == "" {
		opts.RemoteName = "origin"
	}

	if err = tree.Pull(opts); err != nil {
		if !errors.Is(err, gogit.NoErrAlreadyUpToDate) {
			return errors.Wrap(err, "pull and merge new commits")
		}
	}

	return nil
}

// Head возвращает голову
func (r *gitRepository) Head() (Ref, error) {
	r.RLock()
	defer r.RUnlock()

	ref, err := r.repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "get repo head")
	}

	return &gitReference{ref: ref}, nil
}

// Checkout переключает репу на ветку либо каммит по хэшу
func (r *gitRepository) Checkout(target string, options ...Option[gogit.CheckoutOptions]) error {
	r.Lock()
	defer r.Unlock()

	tree, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "getting work tree")
	}

	opts := processOptions(options...)

	if plumbing.IsHash(target) {
		opts.Hash = plumbing.NewHash(target)
	} else {
		opts.Branch = plumbing.NewBranchReferenceName(target)
	}

	if err = tree.Checkout(opts); err != nil {
		return errors.Wrap(err, "checkout to commit")
	}

	return nil
}

func (r *gitRepository) Add(opts ...Option[gogit.AddOptions]) error {
	options := processOptions(opts...)

	r.Lock()
	defer r.Unlock()

	tree, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "get work tree")
	}

	if err = tree.AddWithOptions(options); err != nil {
		return errors.Wrap(err, "add files to tree")
	}

	return nil
}

func (r *gitRepository) Commit(message string, opts ...Option[gogit.CommitOptions]) (string, error) {
	options := processOptions(opts...)

	if options.Author == nil {
		options.Author = &object.Signature{
			Name:  r.config.User.Name,
			Email: r.config.User.Email,
			When:  time.Now(),
		}
	}

	r.Lock()
	defer r.Unlock()

	tree, err := r.repo.Worktree()
	if err != nil {
		return "", errors.Wrap(err, "get work tree")
	}

	var status gogit.Status

	if status, err = tree.Status(); err != nil {
		return "", errors.Ctx().Any("status", status).Wrap(err, "check work tree status")
	}

	var commit plumbing.Hash

	if commit, err = tree.Commit(message, options); err != nil {
		return "", errors.Wrap(err, "send commit")
	}

	return commit.String(), nil
}

func (r *gitRepository) Push(options ...Option[gogit.PushOptions]) error {
	opts := processOptions(options...)
	opts.Auth = r.auth

	r.Lock()
	defer r.Unlock()

	if err := r.repo.Push(opts); err != nil {
		return errors.Wrap(err, "send push")
	}

	return nil
}

// Remotes возвращает ремоуты репозитория из его конфига
func (r *gitRepository) Remotes() map[string][]string {
	r.RLock()
	defer r.RUnlock()

	result := make(map[string][]string)

	if r.config != nil {
		for name, remote := range r.config.Remotes {
			result[name] = remote.URLs
		}
	}

	return result
}

package git

import (
	"time"

	"git.eth4.dev/golibs/errors"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type gitRepository struct {
	path   string
	config *config.Config
	repo   *gogit.Repository
}

// Pull пуллит репу
func (r *gitRepository) Pull() (Repo, error) {
	repo, err := gogit.PlainOpen(r.path)
	if err != nil {
		return nil, errors.Wrap(err, "pull repo")
	}

	r.repo = repo

	r.config, err = repo.Config()
	if err != nil {
		return nil, errors.Wrap(err, "get repo config")
	}

	return r, nil
}

// Head возвращает голову
func (r *gitRepository) Head() (Ref, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "get repo head")
	}

	return &gitReference{ref: ref}, nil
}

// Checkout переключает репу на ветку либо каммит по хэшу
func (r *gitRepository) Checkout(target string, options ...CheckoutOption) error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "getting work tree")
	}

	opts := processOptions[CheckoutOption, *gogit.CheckoutOptions](options...)

	if plumbing.IsHash(target) {
		opts = &gogit.CheckoutOptions{
			Hash: plumbing.NewHash(target),
		}
	} else {
		opts = &gogit.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(target),
		}
	}

	if err = wt.Checkout(opts); err != nil {
		return errors.Wrap(err, "checkout to commit")
	}

	return nil
}

func (r *gitRepository) Add(opts ...AddOption) error {
	options := processOptions[AddOption, *gogit.AddOptions](opts...)

	tree, err := r.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "get work tree")
	}

	if err = tree.AddWithOptions(options); err != nil {
		return errors.Wrap(err, "add files to tree")
	}

	return nil
}

func (r *gitRepository) Commit(message string, opts ...CommitOption) (string, error) {
	options := processOptions[CommitOption, *gogit.CommitOptions](opts...)

	if options.Author == nil {
		options.Author = &object.Signature{
			Name:  r.config.User.Name,
			Email: r.config.User.Email,
			When:  time.Now(),
		}
	}

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

func (r *gitRepository) Push(options ...PushOption) error {
	opts := processOptions[PushOption, *gogit.PushOptions](options...)

	if err := r.repo.Push(opts); err != nil {
		return errors.Wrap(err, "send push")
	}

	return nil
}

// Remotes возвращает ремоуты репозитория из его конфига
func (r *gitRepository) Remotes() map[string][]string {
	result := make(map[string][]string)

	if r.config != nil {
		for name, remote := range r.config.Remotes {
			result[name] = remote.URLs
		}
	}

	return result
}

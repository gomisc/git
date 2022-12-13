package git

import (
	"os"
	"path/filepath"

	"git.eth4.dev/golibs/errors"
	"git.eth4.dev/golibs/filepaths"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type (
	TokenName string

	Ref interface {
		Name() string
		Hash() string
	}

	Repo interface {
		Pull() (Repo, error)
		Head() (Ref, error)
		Checkout(target string, options ...Option[gogit.CheckoutOptions]) error
		Add(opts ...Option[gogit.AddOptions]) error
		Commit(message string, opts ...Option[gogit.CommitOptions]) (string, error)
		Push(options ...Option[gogit.PushOptions]) error
		Remotes() map[string][]string
	}

	Auth interface {
		BasicMethod() *http.BasicAuth
		TokenMethod() *http.TokenAuth
	}
)

// Open - открывает существующий репозиторий
func Open(path string, opts ...Option[gogit.CloneOptions]) (Repo, error) {
	var (
		repo Repo = &gitRepository{path: path}
		err  error
	)

	switch {
	case !filepaths.FileExists(path):
		return Clone(path, opts...)
	case !filepaths.FileExists(filepath.Join(path, ".git")):
		return nil, ErrWrongRepoPath
	}

	repo, err = repo.Pull()
	if err != nil {
		return nil, errors.Wrap(err, "open repo")
	}

	return repo, nil
}

// Clone - клонирует репозиторий с укзанным uri по указаанному пути
// и возвращает интерфейс доступа к нему
func Clone(path string, options ...Option[gogit.CloneOptions]) (Repo, error) {
	if filepaths.FileExists(path) {
		if filepaths.FileExists(filepath.Join(path, ".git")) {
			return nil, ErrRepoPathNotEmpty
		}
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "create repo path")
	}

	opts := processOptions(options...)

	repo, err := gogit.PlainClone(path, false, opts)
	if err != nil {
		return nil, errors.Wrap(err, "clone repository")
	}

	var conf *config.Config

	conf, err = repo.Config()
	if err != nil {
		return nil, errors.Wrap(err, "get repo config")
	}

	return &gitRepository{
		path:   path,
		repo:   repo,
		config: conf,
	}, nil
}

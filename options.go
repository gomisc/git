package git

import (
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type (
	CloneOption    func(o *gogit.CloneOptions)
	CheckoutOption func(o *gogit.CheckoutOptions)
	AddOption      func(o *gogit.AddOptions)
	CommitOption   func(o *gogit.CommitOptions)
	PushOption     func(o *gogit.PushOptions)

	OptionFunc interface {
		CloneOption |
			CheckoutOption |
			AddOption |
			CommitOption |
			PushOption
	}

	OptionType interface {
		*gogit.CloneOptions |
			*gogit.CheckoutOptions |
			*gogit.AddOptions |
			*gogit.CommitOptions |
			*gogit.PushOptions
	}
)

func WithAuth(auth Auth) CloneOption {
	return func(o *gogit.CloneOptions) {
		o.Auth = auth.BasicMethod()
	}
}

func WithURI(uri string) CloneOption {
	return func(o *gogit.CloneOptions) {
		o.URL = uri
	}
}

func WithBranch(branch string) CloneOption {
	return func(o *gogit.CloneOptions) {
		o.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}
}

func WithTag(tag string) CloneOption {
	return func(o *gogit.CloneOptions) {
		o.ReferenceName = plumbing.NewTagReferenceName(tag)
	}
}

func WithCreate() CheckoutOption {
	return func(o *gogit.CheckoutOptions) {
		o.Create = true
	}
}

func ForceCheckout() CheckoutOption {
	return func(o *gogit.CheckoutOptions) {
		o.Force = true
	}
}

func KeepCheckout() CheckoutOption {
	return func(o *gogit.CheckoutOptions) {
		o.Keep = true
	}
}

func AddAll() AddOption {
	return func(o *gogit.AddOptions) {
		o.All = true
	}
}

func AddPath(p string) AddOption {
	return func(o *gogit.AddOptions) {
		o.Path = p
	}
}

func AddMask(m string) AddOption {
	return func(o *gogit.AddOptions) {
		o.Glob = m
	}
}

func WithAuthor(user, email string, t time.Time) CommitOption {
	return func(o *gogit.CommitOptions) {
		o.Author = &object.Signature{
			Name:  user,
			Email: email,
			When:  t,
		}
	}
}

func WithCommiter(user, email string, t time.Time) CommitOption {
	return func(o *gogit.CommitOptions) {
		o.Committer = &object.Signature{
			Name:  user,
			Email: email,
			When:  t,
		}
	}
}

func CommitAll() CommitOption {
	return func(o *gogit.CommitOptions) {
		o.All = true
	}
}

func ForcePush() PushOption {
	return func(o *gogit.PushOptions) {
		o.Force = true
	}
}

func WithRemote(r string) PushOption {
	return func(o *gogit.PushOptions) {
		o.RemoteName = r
	}
}

func WithRemoteURL(rurl string) PushOption {
	return func(o *gogit.PushOptions) {
		o.RemoteURL = rurl
	}
}

func WithTags() PushOption {
	return func(o *gogit.PushOptions) {
		o.FollowTags = true
	}
}

func WithAtomic() PushOption {
	return func(o *gogit.PushOptions) {
		o.Atomic = true
	}
}

func WithOpts(opts map[string]string) PushOption {
	return func(o *gogit.PushOptions) {
		o.Options = opts
	}
}

func processOptions[F OptionFunc, T OptionType](opts ...F) T {
	options := *new(T)

	for i := 0; i < len(opts); i++ {
		opts[i](options)
	}

	return options
}

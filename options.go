package git

import (
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type (
	OptionType interface {
		gogit.CloneOptions |
			gogit.PullOptions |
			gogit.CheckoutOptions |
			gogit.AddOptions |
			gogit.CommitOptions |
			gogit.PushOptions
	}

	Option[T OptionType] func(*T)
)

func CloneAuth(auth Auth) Option[gogit.CloneOptions] {
	return func(o *gogit.CloneOptions) {
		o.Auth = auth.BasicMethod()
	}
}

func CloneURI(uri string) Option[gogit.CloneOptions] {
	return func(o *gogit.CloneOptions) {
		o.URL = uri
	}
}

func CloneBranch(branch string) Option[gogit.CloneOptions] {
	return func(o *gogit.CloneOptions) {
		o.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}
}

func CloneTag(tag string) Option[gogit.CloneOptions] {
	return func(o *gogit.CloneOptions) {
		o.ReferenceName = plumbing.NewTagReferenceName(tag)
	}
}

func PullRemote(name string) Option[gogit.PullOptions] {
	return func(o *gogit.PullOptions) {
		o.RemoteName = name
	}
}

func PullRemoteURL(uri string) Option[gogit.PullOptions] {
	return func(o *gogit.PullOptions) {
		o.RemoteURL = uri
	}
}

func PullBranch(branch string) Option[gogit.PullOptions] {
	return func(o *gogit.PullOptions) {
		o.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}
}

func PullTag(tag string) Option[gogit.PullOptions] {
	return func(o *gogit.PullOptions) {
		o.ReferenceName = plumbing.NewTagReferenceName(tag)
	}
}

func CheckoutCreate() Option[gogit.CheckoutOptions] {
	return func(o *gogit.CheckoutOptions) {
		o.Create = true
	}
}

func ForceCheckout() Option[gogit.CheckoutOptions] {
	return func(o *gogit.CheckoutOptions) {
		o.Force = true
	}
}

func KeepCheckout() Option[gogit.CheckoutOptions] {
	return func(o *gogit.CheckoutOptions) {
		o.Keep = true
	}
}

func AddAll() Option[gogit.AddOptions] {
	return func(o *gogit.AddOptions) {
		o.All = true
	}
}

func AddPath(p string) Option[gogit.AddOptions] {
	return func(o *gogit.AddOptions) {
		o.Path = p
	}
}

func AddMask(m string) Option[gogit.AddOptions] {
	return func(o *gogit.AddOptions) {
		o.Glob = m
	}
}

func CommitAuthor(user, email string, t time.Time) Option[gogit.CommitOptions] {
	return func(o *gogit.CommitOptions) {
		o.Author = &object.Signature{
			Name:  user,
			Email: email,
			When:  t,
		}
	}
}

func Commiter(user, email string, t time.Time) Option[gogit.CommitOptions] {
	return func(o *gogit.CommitOptions) {
		o.Committer = &object.Signature{
			Name:  user,
			Email: email,
			When:  t,
		}
	}
}

func CommitAll() Option[gogit.CommitOptions] {
	return func(o *gogit.CommitOptions) {
		o.All = true
	}
}

func ForcePush() Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.Force = true
	}
}

func PushRemote(r string) Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.RemoteName = r
	}
}

func PushRemoteURL(uri string) Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.RemoteURL = uri
	}
}

func PushTags() Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.FollowTags = true
	}
}

func PushAtomic() Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.Atomic = true
	}
}

func PushOpts(opts map[string]string) Option[gogit.PushOptions] {
	return func(o *gogit.PushOptions) {
		o.Options = opts
	}
}

func processOptions[T OptionType](opts ...Option[T]) *T {
	var options = new(T)

	for _, optFn := range opts {
		optFn(options)
	}

	return options
}

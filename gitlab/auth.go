package gitlab

import (
	"git.eth4.dev/golibs/git"
)

const (
	AccessToken = "access-token"
	DeployToken = "deploy-token"
	JobToken    = "job-token"
)

type gitlabAuth struct {
	username string
	tokens   map[git.TokenName]string
}

func NewAuth(user string, access, deploy, job string) git.Auth {
	return &gitlabAuth{
		username: user,
		tokens: map[git.TokenName]string{
			AccessToken: access,
			DeployToken: deploy,
			JobToken:    job,
		},
	}
}

func (ga *gitlabAuth) Username() string {
	return ga.username
}

func (ga *gitlabAuth) Token(name git.TokenName) string {
	return ga.tokens[name]
}

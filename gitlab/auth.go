package gitlab

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"gopkg.in/gomisc/git.v1"
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

func (ga *gitlabAuth) BasicMethod() *http.BasicAuth {
	return &http.BasicAuth{
		Username: ga.username,
		Password: ga.getToken(),
	}
}

func (ga *gitlabAuth) TokenMethod() *http.TokenAuth {
	return &http.TokenAuth{
		Token: ga.getToken(),
	}
}

func (ga *gitlabAuth) getToken() string {
	switch {
	case ga.tokens[AccessToken] != "":
		return ga.tokens[AccessToken]
	case ga.tokens[JobToken] != "":
		return ga.tokens[JobToken]
	case ga.tokens[DeployToken] != "":
		return ga.tokens[DeployToken]
	default:
		return ""
	}
}

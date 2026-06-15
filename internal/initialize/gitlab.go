package initialize

import (
	"license/internal/gitlab"
)

func InitGitLabCert() error {
	return gitlab.LoadKeys()
}

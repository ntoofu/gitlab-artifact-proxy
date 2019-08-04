package glartpxy

import (
	"github.com/xanzy/go-gitlab"
	"io"
)

type GitlabArtifactAPI interface {
	GetSucceededJobs(project string) ([]gitlab.Job, error)
	DownloadArtifact(project string, ref string, job string) (io.Reader, error)
}

type GitlabArtifactAPIClient struct {
	Host  string
	Token string
}

func (c GitlabArtifactAPIClient) GetSucceededJobs(project string) ([]gitlab.Job, error) {
	// TODO: impliment
	return nil, nil
}

func (c GitlabArtifactAPIClient) DownloadArtifact(project string, ref string, job string) (io.Reader, error) {
	// TODO: impliment
	return nil, nil
}

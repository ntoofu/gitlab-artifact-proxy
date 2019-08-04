package glartpxy

import (
	"errors"
	// "testing"
	"io"
	"os"
	"time"

	"github.com/xanzy/go-gitlab"
)

type GitlabArtifactAPIStub struct {
	stubJobResponse []gitlab.Job
	stubArtifact    string
}

func createGitlabArtifactAPIStub() *GitlabArtifactAPIStub {
	t1 := time.Date(2019, 8, 4, 6, 30, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	stub := GitlabArtifactAPIStub{}
	stub.stubJobResponse = []gitlab.Job{
		gitlab.Job{
			Commit: &gitlab.Commit{
				ID:            "0123456789abcdef0123456789abcdef01234567",
				CommittedDate: &t1,
			},
			Name:   "job1",
			Ref:    "master",
			Status: "success",
		},
		gitlab.Job{
			Commit: &gitlab.Commit{
				ID:            "0101010101010101010101010101010101010101",
				CommittedDate: &t2,
			},
			Name:   "job1",
			Ref:    "master",
			Status: "success",
		},
		gitlab.Job{
			Commit: &gitlab.Commit{
				ID:            "0000111122223333444455556666777788889999",
				CommittedDate: &t3,
			},
			Name:   "job2",
			Ref:    "master",
			Status: "success",
		},
	}
	stub.stubArtifact = "test.zip"
	return &stub
}

func (c GitlabArtifactAPIStub) GetSucceededJobs(project string) ([]gitlab.Job, error) {
	return c.stubJobResponse, nil
}

func (c GitlabArtifactAPIStub) DownloadArtifact(project string, ref string, job string) (io.Reader, error) {
	if project == "some-team%2fsome-project" && ref == "master" && job == "job1" {
		f, err := os.Open(c.stubArtifact)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("Unexpected arguments")
}

package glartpxy

import (
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type ArtifactCache struct {
	TimeToLive    time.Duration
	gitlabClient  GitlabArtifactAPI
	artifact      ArtifactIdentifier
	commitSHA1    string
	verifiedTime  time.Time
	cacheFilePath string
	mutex         *sync.RWMutex
}

type Cache struct {
	os.File
	mutex *sync.RWMutex
}

func (c *Cache) Close() error {
	err := c.File.Close()
	c.mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func NewArtifactCache(ttl time.Duration, gitlabClient GitlabArtifactAPI, artifact ArtifactIdentifier) *ArtifactCache {
	return &ArtifactCache{
		TimeToLive:    ttl,
		gitlabClient:  gitlabClient,
		artifact:      artifact,
		commitSHA1:    "",
		verifiedTime:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		cacheFilePath: "",
		mutex:         &sync.RWMutex{},
	}
}

type fileWithMutex struct {
	os.File
	mutex *sync.RWMutex
}

func (f *fileWithMutex) Close() error {
	err := f.File.Close()
	f.mutex.RUnlock()
	return err
}

func (ac *ArtifactCache) Open() (ReadAtCloser, error) {
	err := ac.updateIfStale()
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred during cache update")
	}
	ac.mutex.RLock()
	f, err := os.Open(ac.cacheFilePath)
	if err != nil {
		ac.mutex.RUnlock()
		return nil, errors.Wrapf(err, "Failed to open cache file '%s'", ac.cacheFilePath)
	}
	return f, nil
}

func (ac *ArtifactCache) updateIfStale() error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	now := time.Now()
	if !now.After(ac.verifiedTime.Add(ac.TimeToLive)) {
		return nil
	}

	jobs, err := ac.gitlabClient.GetSucceededJobs(ac.artifact.Project)
	if err != nil {
		return errors.Wrap(err, "Failed to get succeeded job list")
	}
	commit, err := findLatestCommitOfJob(jobs, ac.artifact.Ref, ac.artifact.Job)
	if err != nil {
		return errors.Wrap(err, "Failed to find target job")
	}
	upstreamCommit := commit.ID
	cacheCommit := ac.commitSHA1
	if upstreamCommit != cacheCommit {
		err = ac.update(upstreamCommit)
		if err != nil {
			return errors.Wrap(err, "Failed to update cache")
		}
	}
	ac.verifiedTime = now
	return nil
}

func (ac *ArtifactCache) update(newCommitSHA1 string) error {
	r, err := ac.gitlabClient.DownloadArtifact(ac.artifact)
	if err != nil {
		return errors.Wrap(err, "Failed to download artifact")
	}

	f, err := ioutil.TempFile("/tmp/", "artifact_")
	if err != nil {
		return errors.Wrap(err, "Failed to create new cache file")
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return errors.Wrap(err, "An error occurred while writing content of artifact to cache file")
	}

	oldCacheFilePath := ac.cacheFilePath
	ac.commitSHA1 = newCommitSHA1
	ac.cacheFilePath = f.Name()
	os.Remove(oldCacheFilePath)
	return nil
}

func findLatestCommitOfJob(jobs []gitlab.Job, ref string, job string) (*gitlab.Commit, error) {
	var commit *gitlab.Commit
	for _, j := range jobs {
		if j.Ref != ref || j.Name != job {
			continue
		}
		if commit == nil || j.Commit.CommittedDate.After(*commit.CommittedDate) {
			commit = j.Commit
		}
	}
	if commit == nil {
		return nil, errors.New("No job found for given name and ref")
	}
	return commit, nil
}

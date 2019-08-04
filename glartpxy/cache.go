package glartpxy

import (
	"os"
	"sync"
	"time"
)

type ArtifactCache struct {
	TimeToLive    time.Duration
	gitlab        GitlabArtifactAPI
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
	// TODO: Error handling
	c.File.Close()
	c.mutex.Unlock()
	return nil
}

func NewArtifactCache(ttl time.Duration, gitlab GitlabArtifactAPI, artifact ArtifactIdentifier) *ArtifactCache {
	return &ArtifactCache{
		TimeToLive:    ttl,
		gitlab:        gitlab,
		artifact:      artifact,
		commitSHA1:    "",
		verifiedTime:  time.Now(),
		cacheFilePath: "",
		mutex:         &sync.RWMutex{},
	}
}

func (ac *ArtifactCache) Open() (ReadAtCloser, error) {
	ac.updateIfStale()
	// R lock
	// get cache file
	// wrap (to unlock) and return
	return nil, nil
}

func (ac *ArtifactCache) updateIfStale() error {
	// TODO: Error handlling
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	// check cache
	//   stale -> update(write new file & remove old one)
	return nil
}

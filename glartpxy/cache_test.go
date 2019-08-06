package glartpxy

import (
	"bytes"
	"github.com/xanzy/go-gitlab"
	"io"
	"os"
	"testing"
	"time"
)

func readeratToBuffer(r io.ReaderAt, b *bytes.Buffer) error {
	buf := make([]byte, 4096)
	var offset int64 = 0
	for {
		n, err := r.ReadAt(buf, offset)
		b.Write(buf[:n])
		if err == io.EOF {
			return nil
		}
		if n < len(buf) {
			return err
		}
		offset += int64(n)
	}
}

func TestOpen(t *testing.T) {
	expected, err := os.Open("test.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer expected.Close()

	gl := createGitlabArtifactAPIStub()
	artifact := ArtifactIdentifier{"some-team%2fsome-project", "master", "job1"}
	cache := NewArtifactCache(time.Second*5, gl, artifact)
	c, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	bufCache := new(bytes.Buffer)
	bufExpected := new(bytes.Buffer)
	readeratToBuffer(c, bufCache)
	readeratToBuffer(expected, bufExpected)
	if bytes.Compare(bufCache.Bytes(), bufExpected.Bytes()) != 0 {
		t.Fatal("Content served by cache is wrong")
	}
}

func TestCacheTTL(t *testing.T) {
	gl := createGitlabArtifactAPIStub()
	artifact := ArtifactIdentifier{"some-team%2fsome-project", "master", "job1"}
	cache := NewArtifactCache(time.Second*5, gl, artifact)
	c1, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	bufCache1 := new(bytes.Buffer)
	readeratToBuffer(c1, bufCache1)
	c1.Close()

	newTime := time.Date(2019, 8, 5, 0, 0, 0, 0, time.UTC)
	gl.stubJobResponse = append(
		gl.stubJobResponse,
		gitlab.Job{
			Commit: &gitlab.Commit{
				ID:            "ffffffffffffffffffffffffffffffffffffffff",
				CommittedDate: &newTime,
			},
			Name:   "job1",
			Ref:    "master",
			Status: "success",
		})
	gl.stubArtifact = "test2.zip"

	// expected not to exceed cache TTL
	time.Sleep(time.Second * 1)

	c2, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	bufCache2 := new(bytes.Buffer)
	readeratToBuffer(c2, bufCache2)
	c2.Close()

	if bytes.Compare(bufCache1.Bytes(), bufCache2.Bytes()) != 0 {
		t.Fatal("Content served by cache has unexpectedly changed though it was not stale")
	}

	// expected to exceed cache TTL
	time.Sleep(time.Second * 5)

	c3, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	bufCache3 := new(bytes.Buffer)
	readeratToBuffer(c3, bufCache3)
	c3.Close()

	if bytes.Compare(bufCache1.Bytes(), bufCache3.Bytes()) == 0 {
		t.Fatal("Content served by cache has not been updated")
	}
}

func TestUpdateCheck(t *testing.T) {
	gl := createGitlabArtifactAPIStub()
	artifact := ArtifactIdentifier{"some-team%2fsome-project", "master", "job1"}
	cache := NewArtifactCache(time.Second*5, gl, artifact)
	c1, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	bufCache1 := new(bytes.Buffer)
	readeratToBuffer(c1, bufCache1)
	c1.Close()

	// to verify codes for checking metadata,
	// change artifact without metadata update
	gl.stubArtifact = "test2.zip"

	// expected to exceed cache TTL
	time.Sleep(time.Second * 2)

	c2, err := cache.Open()
	if err != nil {
		t.Fatal(err)
	}
	bufCache2 := new(bytes.Buffer)
	readeratToBuffer(c2, bufCache2)
	c2.Close()

	if bytes.Compare(bufCache1.Bytes(), bufCache2.Bytes()) != 0 {
		t.Fatal("Content served by cache has unexpectedly changed, though its metadata indicated it was not updated")
	}
}

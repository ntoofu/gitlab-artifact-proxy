package glartpxy

import (
	"bytes"
	"testing"
)

func TestGetFile(t *testing.T) {
	sv := CreateGilabArtifactServer(GitlabArtifactAPIStub{})
	var buf bytes.Buffer
	err := sv.GetFile(ArtifactIdentifier{"some-team%2fsome-project", "master", "job1"}, "dir1/file1", &buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(buf)
}

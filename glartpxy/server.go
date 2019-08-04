package glartpxy

import (
	"io"
	"sync"
)

type GitlabArtifactServer struct {
	GitlabClient  GitlabArtifactAPI
	cacheMetadata *sync.Map // use like map[ArtifactIdentifier]*artifactCache
}

func CreateGilabArtifactServer(client GitlabArtifactAPI) *GitlabArtifactServer {
	return &GitlabArtifactServer{client, &sync.Map{}}
}

func (sv GitlabArtifactServer) GetFile(artifact ArtifactIdentifier, filepath string, w io.Writer) error {
	// find cache
	//   no entry -> create new
	// getcache
	// defer close cache
	// unzip & write
	// return
	return nil
}

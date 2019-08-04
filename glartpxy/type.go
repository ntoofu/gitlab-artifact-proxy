package glartpxy

import (
	"io"
)

type ReadAtCloser interface {
	io.ReaderAt
	io.Closer
}

type ArtifactIdentifier struct {
	Project string
	Ref     string
	Job     string
}

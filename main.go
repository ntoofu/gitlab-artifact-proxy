package main

import (
	"github.com/ntoofu/gitlab-artifact-proxy/glartpxy"
	"io"
	"net/http"
)

var server glartpxy.GitlabArtifactServer

func gitlabArtifactProxy(w http.ResponseWriter, req *http.Request) {
	// TODO: parse path in request
	project := ""
	ref := ""
	job := ""
	filepath := ""

	f, err := server.GetFile(project, ref, job, filepath, w)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	return
}

func main() {
	server := glartpxy.GitlabArtifactServer(glartpxy.GitlabArtifactAPIClient("foo.bar.com", "TokenXXX"))
	http.HandleFunc("/", gitlabArtifactProxy)
	http.ListenAndServe(":8080", nil)
}

package glartpxy

import (
	"bytes"
	"os"
	"testing"
)

func TestOpenFileInZipArchive(t *testing.T) {
	zip, err := os.Open("test.zip")
	buf := make([]byte, 32)
	if err != nil {
		t.Fatal(err)
	}
	f, err := OpenFileInZipArchive(zip, "dir1/file1")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.Read(buf)
	if bytes.Compare(buf, []byte("path: dir1>file1")) != 0 {
		t.Fatalf("Content in the unarchived file is wrong")
	}
}

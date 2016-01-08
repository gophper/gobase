package gobase

import (
	"testing"
)

func TestFilePutContent(t *testing.T) {
	FilePutContent("/tmp/cs.txt", "cs")

	fg, _ := FileGetContent("/tmp/cs.txt")

	if fg != "cs" {
		t.Error("FilePutContent and FileGetContent is Error")
	}

	Remove("/tmp/cs.txt")

	if IsFile("/tmp/cs.txt") {
		t.Error("Remove file Error")
	}

}

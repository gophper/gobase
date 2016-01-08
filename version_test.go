package gobase

import (
	"testing"
)

func TestVersion(t *testing.T) {

	if Version() != "gobase 1.0" {
		t.Error("Func Version is Error")
	}

}

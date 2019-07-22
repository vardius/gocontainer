package gocontainer

import (
	"testing"
)

func TestNew(t *testing.T) {
	container := New()

	if container == nil {
		t.Fail()
	}
}

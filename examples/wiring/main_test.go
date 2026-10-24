package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/vardius/gocontainer/v2"
)

func TestRegister(t *testing.T) {
	var c gocontainer.Container
	var logs strings.Builder
	register(&c, &logs)

	greeter := greeterKey.MustGet(&c)
	got := []string{greeter.Greet(1), greeter.Greet(2), greeter.Greet(3)}

	want := []string{"Hello, Ada.", "Hello, Grace.", "Hello, stranger."}
	if !slices.Equal(got, want) {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
	if !strings.Contains(logs.String(), `level=WARN msg="unknown user" id=3`) {
		t.Errorf("logs = %q, want a warning about user 3", logs.String())
	}
}

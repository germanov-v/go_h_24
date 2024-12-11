package main

import (
	"github.com/germanov-v/go_h_24/hw01_hello_otus/stringutil"
	"testing"
)

func TestRevers(t *testing.T) {
	got := stringutil.Reverse("test")
	expected := "tset"

	if got != expected {
		t.Fail()
	}
}

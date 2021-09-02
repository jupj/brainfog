package main

import (
	_ "embed"
	"testing"
)

//go:embed hello.bf
var helloProg []byte

func Test_Helloworld(t *testing.T) {
	bf := newBrainfog(helloProg)
	var res []byte
	for c := range bf.outCh {
		res = append(res, c)
	}

	wantedRes := "Hello World!\n"

	if string(res) != wantedRes {
		t.Errorf("Got %q, want %q", string(res), wantedRes)
	}
}

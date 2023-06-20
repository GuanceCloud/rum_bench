package main

import (
	"fmt"
	"io"
	"testing"
)

func Test(t *testing.T) {
	ct, r := getReplayBody()

	fmt.Println(ct)
	text, _ := io.ReadAll(r)

	fmt.Println(string(text))
}

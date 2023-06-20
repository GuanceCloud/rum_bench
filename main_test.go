package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	txt, err := os.ReadFile("./testdata/segment")
	if err != nil {
		t.Fatal(err)
	}

	str := base64.StdEncoding.EncodeToString(txt)
	fmt.Println("-------------")
	fmt.Println(str)
	fmt.Println("-------------")
}

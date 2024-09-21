package test

import (
	"blog/util"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	digest := util.Md5("123456")
	if digest != "e10adc3949ba59abbe56e057f20f883e" {
		fmt.Println(digest)
		t.Fail()
	}
}

// go test -v .\util\test\ -run=^TestHash$ -count=1

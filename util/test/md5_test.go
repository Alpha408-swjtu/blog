package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	md5 := md5.New()

	text := "你好"
	md5.Write([]byte(text))
	digest := md5.Sum(nil) //md5哈希的结果是128bit，即16个字节
	if len(digest) != 16 {
		t.Fail()
	}
	fmt.Println(digest)

	md5.Reset()
	text = "海上升明月"
	md5.Write([]byte(text))
	digest = md5.Sum(nil)
	if len(digest) != 16 {
		t.Fail()
	}
	fmt.Println(digest)
}

func TestHexEncode(t *testing.T) {
	digest := []byte{126, 202, 104, 159, 13, 51, 137, 217, 222, 166, 106, 225, 18, 229, 207, 215}
	s := hex.EncodeToString(digest) //十六进制编码，1个byte对应到2个字符。[0-9a-f]
	fmt.Println(s)
	if len(s) != 2*len(digest) {
		t.Fail()
	}
	b, _ := hex.DecodeString(s)
	fmt.Println(b)
}

// go test -v .\util\test\ -run=^TestMd5$ -count=1
// go test -v .\util\test\ -run=^TestHexEncode$ -count=1

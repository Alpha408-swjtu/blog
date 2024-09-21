package test

import (
	"blog/util"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/sonic"
)

func TestCamel2Snake(t *testing.T) {
	fmt.Println("hello")
	s1 := "Abc"
	s2 := util.Camel2Snake(s1)
	if s2 != "abc" {
		fmt.Println(s2)
		t.Fail()
	}

	s1 = "AbcEfg"
	s2 = util.Camel2Snake(s1)
	if s2 != "abc_efg" {
		fmt.Println(s2)
		t.Fail()
	}

	s1 = "abcEfg"
	s2 = util.Camel2Snake(s1)
	if s2 != "abc_efg" {
		fmt.Println(s2)
		t.Fail()
	}
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		util.Camel2Snake("UserName")
	}
}

type User struct {
	Name     string
	Age      int
	height   float32
	Birthday time.Time
}

var user = User{Name: "大乔乔", Age: 18, height: 170.5, Birthday: time.Now()}

func BenchmarkStdJson(b *testing.B) {
	b.ResetTimer() //把本行之前的代码耗时排除在外
	for i := 0; i < b.N; i++ {
		bs, _ := json.Marshal(user)
		var inst User
		json.Unmarshal(bs, &inst)
	}
}

func BenchmarkSonicJson(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs, _ := sonic.Marshal(user)
		var inst User
		sonic.Unmarshal(bs, &inst)
	}
}

func TestRandString(t *testing.T) {
	s := util.RandStringRunes(30)
	fmt.Println(s)
}

func TestStr(t *testing.T) {
	util.Str()
}

// go test -v ./util/test -run=^TestRandString$ -count=1
// go test -v ./util/test -run=^TestStr$ -count=1
// go test -v ./util/test -run=^TestCamel2Snake$ -count=1
// go test -v ./util/test -run=^TestCamel2Snake$ -count=1 -bench=RandString
// go test ./util/test -run=^$ -count=1 -benchmem -bench=Json -benchtime=3s -memprofile=mem -cpuprofile=cpu
// go tool pprof cpu
// go tool pprof mem

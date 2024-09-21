package test

import (
	"blog/database"
	"fmt"
	"testing"
)

func TestRedisToken(t *testing.T) {
	refreshToken, authToken := "123", "abc"
	database.SetToken(refreshToken, authToken)
	token := database.GetToken(refreshToken)
	if token != authToken {
		fmt.Println(token)
		t.Fail()
	}
}

// go test -v .\database\test\ -run=^TestRedisToken$ -count=1

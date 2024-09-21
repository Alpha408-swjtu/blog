package test

import (
	"blog/database"
	"sync"
	"testing"
)

// 测试高并发下gorm.DB是否还是单例
func TestGetBlogDBConnection(t *testing.T) {
	const C = 100
	wg := sync.WaitGroup{}
	wg.Add(C)
	for i := 0; i < C; i++ {
		go func() {
			defer wg.Done()
			database.GetBlogDBConnection()
		}()
	}
	wg.Wait()
}

// go test -v .\database\test\ -run=^TestGetBlogDBConnection$ -count=1

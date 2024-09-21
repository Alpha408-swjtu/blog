package test

import (
	"blog/database"
	"fmt"
	"testing"
)

func TestGetBlogById(t *testing.T) {
	id := 1
	blog := database.GetBlogById(id)
	if blog == nil {
		t.Fail()
	} else {
		fmt.Printf("%+v\n", *blog)
	}
}

func TestGetBlogByUserId(t *testing.T) {
	uid := 1
	blogs := database.GetBlogByUserId(uid)
	if len(blogs) == 0 {
		t.Fail()
	} else {
		for _, blog := range blogs {
			fmt.Println(blog.Id, blog.Title)
		}
	}
}

func TestUpdateBlog(t *testing.T) {
	blog := database.Blog{Id: 1, Title: "双十一", Article: "双十一来临喜洋洋，购物狂欢乐无边。电商盛宴满眼芳，心愿成真喜笑颜。"}
	if err := database.UpdateBlog(&blog); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

// go test -v .\database\test\ -run=^TestGetBlogById$ -count=1
// go test -v .\database\test\ -run=^TestGetBlogByUserId$ -count=1
// go test -v .\database\test\ -run=^TestUpdateBlog$ -count=1

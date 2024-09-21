package handler

import (
	"blog/database"
	"blog/handler/middleware"
	"blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取用户的博客列表
func BlogList(ctx *gin.Context) {
	uid, err := strconv.Atoi(ctx.Param("uid")) //获取restful参数
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid uid")
		return
	}
	blogs := database.GetBlogByUserId(uid)
	util.LogRus.Debugf("get %d blogs of user %d", len(blogs), uid)
	ctx.HTML(http.StatusOK, "blog_list.html", blogs) //go template
}

// 获取某一篇博客的详情
func BlogDetail(ctx *gin.Context) {
	blogId := ctx.Param("bid") //获取restful参数
	bid, err := strconv.Atoi(blogId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid blog id")
		return
	}
	blog := database.GetBlogById(bid)
	if blog == nil {
		ctx.String(http.StatusNotFound, "博客不存在")
		return
	}
	util.LogRus.Debug(blog.Article)
	ctx.HTML(http.StatusOK, "blog.html", gin.H{"title": blog.Title, "article": blog.Article, "bid": blogId, "update_time": blog.UpdateTime.Format("2006-01-02 15:04:05")})
}

// 参数向结构体映射，并执行校验
type UpdateRequest struct {
	BlogId  int    `form:"bid" binding:"gt=0"`     //数字的值大于0
	Title   string `form:"title" binding:"gt=0"`   //字符串长度大于0
	Article string `form:"article" binding:"gt=0"` //字符串长度大于0
}

// 更新博客
func BlogUpdate(ctx *gin.Context) {
	// blogId := ctx.PostForm("bid") //获取post form参数
	// title := ctx.PostForm("title")
	// article := ctx.PostForm("article")
	// bid, err := strconv.Atoi(blogId)
	// if err != nil {
	// 	ctx.String(http.StatusBadRequest, "invalid blog id")
	// 	return
	// }

	var request UpdateRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid parameter")
		return
	}
	bid := request.BlogId
	title := request.Title
	article := request.Article

	blog := database.GetBlogById(bid)
	if blog == nil {
		ctx.String(http.StatusBadRequest, "blog id not exists")
		return
	}
	loginUid := ctx.Value("uid") //从ctx中取得当前登录用户的uid
	if loginUid == nil || loginUid.(int) != blog.UserId {
		ctx.String(http.StatusForbidden, "无权修改")
		return
	}
	err = database.UpdateBlog(&database.Blog{Id: bid, Title: title, Article: article})
	if err != nil {
		util.LogRus.Errorf("update blog %d failed: %s", bid, err)
		ctx.String(http.StatusInternalServerError, "更新失败") //不要把原始的err返回给前端，否则用户通过查看页面源码能看到mysql表结构信息
		return
	}
	ctx.String(http.StatusOK, "success")
}

// 从jwt里解析出uid，判断blog_id是否属于uid
func BlogBelong(ctx *gin.Context) {
	bids := ctx.Query("bid")
	token := ctx.Query("token")
	bid, err := strconv.Atoi(bids)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid blog id")
		return
	}
	blog := database.GetBlogById(bid)
	if blog == nil {
		ctx.String(http.StatusBadRequest, "blog id not exists")
		return
	}
	loginUid := middleware.GetUidFromJwt(token)
	if loginUid == blog.UserId {
		ctx.String(http.StatusOK, "true")
	} else {
		ctx.String(http.StatusOK, "false")
	}
}

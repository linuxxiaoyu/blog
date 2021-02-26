package main

import (
	"context"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	pb "github.com/linuxxiaoyu/blog/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	articleClient pb.ArticleClient
)

func init() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("dial account service failed: %v", err)
	}
	// defer conn.Close()
	articleClient = pb.NewArticleClient(conn)
}

// Articles return 10 articles order by id
// GET /articles
func getArticles(c *gin.Context) {
	req := pb.GetArticlesRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := articleClient.GetArticles(ctx, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// articles := make([]gin.H, 10)
	// for _, article := range r.Rs {
	// 	articles = append(articles, gin.H{
	// 		"uid":     article.GetUid(),
	// 		"title":   article.GetTitle(),
	// 		"context": article.GetContent(),
	// 		"time":    article.GetTime(),
	// 	})
	// }
	c.JSON(http.StatusOK, r.Rs)
}

// Get an article
// GET /articles/:id
func getArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	req := pb.GetArticleRequest{Id: uint32(id)}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := articleClient.GetArticle(ctx, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title":   r.GetTitle(),
		"time":    r.GetTime(),
		"content": r.GetContent(),
		"uid":     r.GetUid(),
	})
}

// New create an article.
// POST /article
// form: token title content
func newArticle(c *gin.Context) {
	title := html.EscapeString(c.PostForm("title"))
	content := html.EscapeString(c.PostForm("content"))

	valid := validation.Validation{}
	valid.Required(title, "title").Message("title can't empty")
	valid.Required(content, "content").Message("content can't empty")

	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	uid, _ := c.Get("uid")
	req := pb.PostArticleRequest{
		Uid:     uid.(uint32),
		Title:   title,
		Content: content,
		Time:    timestamppb.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := articleClient.PostArticle(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": r.GetId(),
	})
}

// Delete an article
// DELETE /articles/:id
// form: token
func deleteArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	req := pb.DeleteArticleRequest{
		Id:  uint32(id),
		Uid: uid.(uint32),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := articleClient.DeleteArticle(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Update an article
// PUT /articles/:id
// form: token title content
func updateArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")
	title := html.EscapeString(c.PostForm("title"))
	content := html.EscapeString(c.PostForm("content"))

	req := pb.PutArticleRequest{
		Id:      uint32(id),
		Uid:     uid.(uint32),
		Title:   title,
		Content: content,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := articleClient.PutArticle(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

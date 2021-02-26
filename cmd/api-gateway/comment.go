package main

import (
	"context"
	"fmt"
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
	commentClient pb.CommentClient
)

func init() {
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("dial comment service failed: %v", err)
	}
	// defer conn.Close()
	commentClient = pb.NewCommentClient(conn)
}

// New a comment by article_id and user_id
// POST /comments
// form: token aid content
func newComment(c *gin.Context) {
	aid, _ := strconv.ParseUint(c.PostForm("aid"), 10, 32)
	uid, _ := c.Get("uid")
	content := html.EscapeString(c.PostForm("content"))

	valid := validation.Validation{}
	valid.Required(content, "content").Message("content can't empty")
	if valid.HasErrors() || aid == 0 {
		fmt.Println(403)
		c.JSON(http.StatusForbidden, nil)
		return
	}

	req := pb.PostCommentRequest{
		Aid:     uint32(aid),
		Uid:     uid.(uint32),
		Content: content,
		Time:    timestamppb.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := commentClient.PostComment(ctx, &req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": r.GetId(),
	})
}

// Delete a comment by comment_id
// DELETE /comments/:id
// form: token
func deleteComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req := pb.DeleteCommentRequest{
		Id:  uint32(id),
		Uid: uid.(uint32),
	}
	_, err := commentClient.DeleteComment(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

// Update a comment by comment_id
// PUT /comments/:id
// form: token content
func updateComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	content := html.EscapeString(c.PostForm("content"))
	valid := validation.Validation{}
	valid.Required(content, "content").Message("content can't empty")
	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	req := pb.PutCommentRequest{
		Id:      uint32(id),
		Uid:     uid.(uint32),
		Content: content,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := commentClient.PutComment(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

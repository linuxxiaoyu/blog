package main

import (
	"context"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	pb "github.com/linuxxiaoyu/blog/api"
	"google.golang.org/grpc"
)

var (
	client pb.AccountClient
)

func init() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("dial account service failed: %v", err)
	}
	// defer conn.Close()
	client = pb.NewAccountClient(conn)
}

// Login a user, return a token
// GET /user
// query: name password
func Login(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}

	valid.Required(name, "name").Message("name can't empty")
	valid.MaxSize(name, 16, "name").Message("name's max length is 16")

	valid.Required(password, "password").Message("password can't empty")
	valid.MinSize(password, 6, "password").Message("password's min length is 6")
	valid.MaxSize(password, 16, "password").Message("password's max length is 16")

	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	req := pb.AccountRequest{
		Name:     name,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetToken(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": r.GetToken(),
	})
}

// Register a user
// POST /user
// form: name password
func Register(c *gin.Context) {
	name := html.EscapeString(c.PostForm("name"))
	password := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("name can't empty")
	valid.MaxSize(name, 16, "name").Message("name's max length is 16")

	valid.Required(password, "password").Message("password can't empty")
	valid.MinSize(password, 6, "password").Message("password's min length is 6")
	valid.MaxSize(password, 16, "password").Message("password's max length is 16")

	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	var req pb.AccountRequest
	req.Name = name
	req.Password = password

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Register(ctx, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": r.GetId(),
	})
}

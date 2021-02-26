package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	pb "github.com/linuxxiaoyu/blog/api"
)

var client pb.AccountClient

func init() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("dial account service failed: %v", err)
	}
	// defer conn.Close()
	client = pb.NewAccountClient(conn)
}

func forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, nil)
	c.Abort()
}

// form: token
func Auth(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	req := pb.ParseTokenRequest{
		Token: token,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := client.ParseToken(ctx, &req)
	if err != nil {
		fmt.Println(err)
		forbidden(c)
		return
	}

	c.Set("uid", r.GetUid())

	c.Next()
}

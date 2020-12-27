package jwt

import (
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
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial account service failed: %v", err)
	}
	defer conn.Close()
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

	resp := pb.AccountRequest{
		Token: token,
	}
	ctx, cancel := context.WithTimeout(c, 500*time.Second)
	defer cancel()
	r, err := client.ParseToken(ctx, &resp)
	if err != nil {
		forbidden(c)
		return
	}

	c.Set("uid", r.GetId())

	c.Next()
}

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

// getArticles return 10 articles order by id
// GET /articles
func getArticles(c *gin.Context) {
	req := pb.GetArticlesRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := articleClient.GetArticles(ctx, &req)
	if err != nil || len(r.Rs) == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// m[aid] = article info
	m := map[uint32]gin.H{}
	resps := r.GetRs()
	uids := make([]uint32, 0, len(r.Rs))
	aids := make([]uint32, 0, len(r.Rs))
	for _, resp := range resps {
		// append article author uids
		uids = append(uids, resp.GetUid())
		aids = append(aids, resp.GetId())
		m[resp.GetId()] = gin.H{
			"title":   resp.GetTitle(),
			"content": resp.GetContent(),
			"time": resp.GetTime().
				AsTime().Local().Format("2006-01-02 15:04:05"),
		}
	}

	names := map[uint32]string{}
	rC, err := commentClient.GetCommentsByAids(ctx,
		&pb.GetCommentsByAidsRequest{Aids: aids})
	if err == nil {
		for _, comm := range rC.GetComments() {
			infos := comm.GetInfos()
			for _, info := range infos {
				// append comment uids
				uids = append(uids, info.Uid)
			}
		}
	}

	rA, err := accountClient.GetAccounts(ctx,
		&pb.GetAccountsRequest{Ids: uids})
	if err != nil {
		// do nothing
	}
	names = rA.GetNames()
	if err == nil {
		for aid, comm := range rC.GetComments() {
			infos := comm.GetInfos()
			comms := make([]gin.H, 0, len(infos))
			for _, info := range infos {
				comms = append(comms, gin.H{
					"content": info.GetContent(),
					"time":    info.GetTime().AsTime().Local().Format("2006-01-02 15:04:05"),
					"name":    names[info.GetUid()],
				})
			}
			m[aid]["comments"] = comms
		}
	}

	for _, resp := range resps {
		m[resp.GetId()]["author"] = names[resp.Uid]
	}

	articles := []gin.H{}
	for _, article := range m {
		articles = append(articles, article)
	}
	c.JSON(http.StatusOK, articles)
}

// getArticle get an article by id
// GET /articles/:id
func getArticle(c *gin.Context) {
	idStr, b := c.Params.Get("id")
	if !b {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	id64, _ := strconv.ParseUint(idStr, 10, 32)
	if id64 <= 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	id := uint32(id64)
	req := pb.GetArticleRequest{Id: id}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	r, err := articleClient.GetArticle(ctx, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	h := gin.H{}
	h["title"] = r.GetTitle()
	h["content"] = r.GetContent()
	h["time"] = r.GetTime().AsTime().Local().Format("2006-01-02 15:04:05")

	uid := r.GetUid()
	// add comments
	rC, err := commentClient.GetCommentsByAids(ctx,
		&pb.GetCommentsByAidsRequest{
			Aids: []uint32{id},
		})
	if err != nil {
		rA, err := accountClient.GetAccount(ctx, &pb.GetAccountRequest{Id: uid})
		if err == nil {
			h["author"] = rA.GetName()
		}
		c.JSON(http.StatusOK, h)
		return
	}

	infos := rC.GetComments()[id].GetInfos()
	uids := make([]uint32, 0, len(infos)+1)
	// Get all user names, include article author's name.
	uids = append(uids, uid)
	for _, comment := range infos {
		uids = append(uids, comment.GetUid())
	}
	resp, err := accountClient.GetAccounts(ctx, &pb.GetAccountsRequest{Ids: uids})
	if err == nil {
		comms := make([]gin.H, 0, len(infos))
		for _, comment := range infos {
			comms = append(comms, gin.H{
				"name":    resp.Names[comment.GetUid()],
				"content": comment.GetContent(),
				"time":    comment.GetTime().AsTime().Local().Format("2006-01-02 15:04:05"),
			})
		}
		h["comments"] = comms
	}
	h["author"] = resp.Names[uid]

	c.JSON(http.StatusOK, h)
}

// newArticle create an article.
// POST /articles
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

	c.JSON(http.StatusOK, r)
}

// deleteArticle delete an article
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

// updateArticle update an article
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

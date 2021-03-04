package main

import (
	"context"
	"errors"

	pb "github.com/linuxxiaoyu/blog/api"
	"github.com/linuxxiaoyu/blog/internal/data"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct{}

func (s *service) GetArticle(ctx context.Context,
	req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	var resp pb.GetArticleResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	article, err := data.GetArticle(ctx, req.GetId())
	if err != nil {
		return &resp, err
	}
	resp.Uid = article.AuthorID
	resp.Title = article.Title
	resp.Content = article.Content
	resp.Time = timestamppb.New(article.Time)
	return &resp, nil
}

func (s *service) GetArticles(ctx context.Context,
	req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error) {
	var resp pb.GetArticlesResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	articles, err := data.GetArticles(ctx)
	if err != nil {
		return &resp, err
	}
	for _, article := range articles {
		resp.Rs = append(resp.Rs, &pb.GetArticleResponse{
			Id:      article.ID,
			Uid:     article.AuthorID,
			Title:   article.Title,
			Content: article.Content,
			Time:    timestamppb.New(article.Time),
		})
	}
	return &resp, nil
}

func (s *service) PostArticle(ctx context.Context,
	req *pb.PostArticleRequest) (*pb.PostArticleResponse, error) {
	var resp pb.PostArticleResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	id, err := data.CreateArticle(
		ctx,
		req.GetUid(),
		req.GetTitle(),
		req.GetContent(),
		req.GetTime().AsTime(),
	)
	if err != nil {
		return &resp, err
	}
	resp.Id = id
	return &resp, nil
}

func (s *service) DeleteArticle(ctx context.Context,
	req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	var resp pb.DeleteArticleResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	return &resp, data.DeleteArticle(ctx, req.GetId(), req.GetUid())
}

func (s *service) PutArticle(ctx context.Context,
	req *pb.PutArticleRequest) (*pb.PutArticleResponse, error) {
	var resp pb.PutArticleResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	return &resp, data.UpdateArticle(
		ctx,
		req.GetId(),
		req.GetUid(),
		req.GetTitle(),
		req.GetContent(),
	)
}

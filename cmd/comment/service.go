package main

import (
	"context"
	"errors"

	pb "github.com/linuxxiaoyu/blog/api"
	"github.com/linuxxiaoyu/blog/internal/data"
)

type service struct{}

// TODO
func (s *service) GetComment(ctx context.Context,
	req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	return nil, nil
}

func (s *service) PostComment(ctx context.Context,
	req *pb.PostCommentRequest) (*pb.PostCommentResponse, error) {
	var resp pb.PostCommentResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	comment := data.Comment{
		AID:     req.GetAid(),
		UID:     req.GetUid(),
		Content: req.GetContent(),
		Time:    req.GetTime().AsTime(),
	}

	id, err := data.CreateComment(ctx, &comment)
	resp.Id = id
	return &resp, err
}

func (s *service) PutComment(ctx context.Context,
	req *pb.PutCommentRequest) (*pb.PutCommentResponse, error) {
	var resp pb.PutCommentResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	comment := data.Comment{
		ID:      req.GetId(),
		UID:     req.GetUid(),
		Content: req.Content,
	}
	err := data.UpdateComment(ctx, &comment)
	return &resp, err
}

func (s *service) DeleteComment(ctx context.Context,
	req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	var resp pb.DeleteCommentResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	err := data.DeleteComment(ctx, req.GetId(), req.GetUid())
	return &resp, err
}

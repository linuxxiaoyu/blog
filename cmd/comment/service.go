package main

import (
	"context"
	"errors"

	pb "github.com/linuxxiaoyu/blog/api"
	"github.com/linuxxiaoyu/blog/internal/data"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct{}

func (s *service) GetCommentsByAids(ctx context.Context,
	req *pb.GetCommentsByAidsRequest) (*pb.GetCommentsByAidsResponse, error) {
	var resp pb.GetCommentsByAidsResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	comments, err := data.GetCommentsByAids(ctx, req.GetAids())
	if err != nil {
		return &resp, err
	}

	if resp.Comments == nil {
		resp.Comments = make(map[uint32]*pb.CommentInfos)
	}
	for key, values := range comments {
		var infos pb.CommentInfos
		for _, v := range values {
			info := pb.CommentInfo{
				Id:      v.ID,
				Uid:     v.UID,
				Aid:     v.AID,
				Content: v.Content,
				Time:    timestamppb.New(v.Time),
			}
			if infos.Infos == nil {
				infos.Infos = []*pb.CommentInfo{}
			}
			infos.Infos = append(infos.Infos, &info)
		}
		resp.Comments[key] = &infos
	}
	return &resp, nil
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

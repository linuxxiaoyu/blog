package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"time"

	pb "github.com/linuxxiaoyu/blog/api"
	"github.com/linuxxiaoyu/blog/internal/data"
	"golang.org/x/crypto/scrypt"
)

var salt = []byte{0x11, 0x22, 0x33, 0x44, 0x33, 0x77, 0x0d, 0x0a}

type service struct{}

func getUser(id uint32, name string) (user data.User, err error) {
	switch {
	case id > 0:
		user, err = data.GetUser(id)
	case name != "":
		user, err = data.GetUserByName(name)
	}
	return
}

func (s *service) GetAccount(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	var resp pb.AccountResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	user, err := getUser(req.Id, req.Name)
	resp.Id = user.ID
	resp.Name = user.Name
	return &resp, err
}

func (s *service) GetToken(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	var resp pb.AccountResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	user, err := getUser(resp.Id, req.Name)
	if err != nil {
		return &resp, err
	}

	if user.Password != addSalt(req.Password) {
		return &resp, errors.New("auth error")
	}

	token, err := generateToken(user.ID, user.Name)
	if err != nil {
		return &resp, err
	}

	resp.Token = token
	return &resp, err
}

func (s *service) ParseToken(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	var resp pb.AccountResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}
	if req.Token == "" {
		return &resp, errors.New("token is empty")
	}
	claims, err := parseToken(req.Token)
	if err != nil || claims.ExpiresAt < time.Now().Unix() {
		return &resp, errors.New("forbidden")
	}
	user, err := getUser(claims.ID, "")
	if err != nil || user.Name != claims.Name {
		return &resp, errors.New("forbidden")
	}

	resp.Id = claims.ID
	resp.Name = claims.Name
	return &resp, nil
}

func (s *service) Register(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	var resp pb.AccountResponse
	if req == nil {
		return &resp, errors.New("req is null")
	}

	_, err := data.GetUserByName(req.Name)
	if err == nil {
		return &resp, errors.New("user name exist")
	}

	if err != nil && data.IsRowNotFound(err) {
		var user data.User
		user.Name = req.Name
		user.Password = addSalt(req.Password)
		id, err := data.CreateUser(user)
		resp.Id = id
		return &resp, err
	}

	return &resp, err
}

func addSalt(password string) string {
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}

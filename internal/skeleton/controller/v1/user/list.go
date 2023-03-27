// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package user

import (
	"time"

	"github.com/changaolee/skeleton/internal/pkg/log"
	pb "github.com/changaolee/skeleton/pkg/proto/skeleton/v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// List 返回用户列表，只有 root 用户才能获取用户列表
func (ctrl *UserController) List(ctx *gin.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.C(ctx).Infow("List user function called")

	resp, err := ctrl.b.Users().List(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}

	users := make([]*pb.UserInfo, 0, len(resp.Users))
	for _, u := range resp.Users {
		createdAt, _ := time.Parse("2006-01-02 15:04:05", u.CreatedAt)
		updatedAt, _ := time.Parse("2006-01-02 15:04:05", u.UpdatedAt)
		users = append(users, &pb.UserInfo{
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	return &pb.ListUserResponse{TotalCount: resp.TotalCount, Users: users}, nil
}

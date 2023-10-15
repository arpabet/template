/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package api

import (
	"context"
	"github.com/keyvalstore/store"
	"github.com/codeallergy/glue"
	"github.com/sprintframework/template/pkg/pb"
	"reflect"
)


var UserServiceClass = reflect.TypeOf((*UserService)(nil)).Elem()

type UserService interface {
	glue.InitializingBean

	GenerateUserId(ctx context.Context) (string, error)

	CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.UserEntity, error)

	ResetPassword(ctx context.Context, userId string, newPassword string) (email string, err error)

	AuthenticateUser(ctx context.Context, login, password string) (*pb.UserEntity, error)

	GetUser(ctx context.Context, userId string) (*pb.UserEntity, error)

	GetUserIdByLogin(ctx context.Context, login string) (string, error)

	GetUserIdByUsername(ctx context.Context, username string) (string, error)

	GetUserIdByEmail(ctx context.Context, email string) (string, error)

	SaveUser(ctx context.Context, user *pb.UserEntity) error

	RemoveUser(ctx context.Context, userId string) error

	DropUserContent(ctx context.Context, userId string) error

	DoWithUser(ctx context.Context, userId string, cb func(user *pb.UserEntity) error) error

	DumpUser(ctx context.Context, userId string, cb func(entry *store.RawEntry) bool) error

	EnumUsers(ctx context.Context, cb func(user *pb.UserEntity) bool) error

	SaveRecoverCode(ctx context.Context, login string, rc *pb.RecoverCodeEntity, ttlSeconds int) error

	ValidateRecoverCode(ctx context.Context, login string, code string) error
}

var SecurityLogServiceClass = reflect.TypeOf((*SecurityLogService)(nil)).Elem()

type SecurityLogService interface {

	LogEvent(ctx context.Context, userId, eventName, remoteIP, userAgent string) error

	EnumEvents(ctx context.Context, userId string, cb func(item *pb.SecurityLogEntity) bool) error

}

var PageServiceClass = reflect.TypeOf((*PageService)(nil)).Elem()

type PageService interface {

	// ErrPageNotFound on error
	GetPage(ctx context.Context, name string) (*pb.PageEntity, error)

	CreatePage(ctx context.Context, page *pb.AdminPage) error

	UpdatePage(ctx context.Context, page *pb.AdminPage) error

	RemovePage(ctx context.Context, name string) error

	EnumPages(ctx context.Context, cb func(page *pb.PageEntity) bool) error

}

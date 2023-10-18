/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package server

import (
	"context"
	"fmt"
	"github.com/sprintframework/sprintframework/sprintutils"
	"github.com/sprintframework/template/pkg/pb"
	"github.com/sprintframework/template/pkg/service"
	"github.com/sprintframework/sprint"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"strconv"
	"time"
)

func (t *implUIGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {

	entity, err := t.UserService.AuthenticateUser(ctx, req.Login, req.Password)
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err == service.ErrUserInvalidPassword {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "Login", entity.UserId)
		}

	}()

	if err != nil {
		return nil, err
	}

	roles := make(map[string]bool)
	roles["WEB_USER"] = true
	if entity.Role == pb.UserRole_ADMIN {
		roles["WEB_ADMIN"] = true
	}

	token, err := t.AuthorizationMiddleware.GenerateToken(&sprint.AuthorizedUser{
		Username:  entity.Username,
		Roles:     roles,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(t.AccessTokenMinutes)).Unix(),
	})

	if err != nil {
		return nil, err
	}

	refreshToken, err := t.AuthorizationMiddleware.GenerateToken(&sprint.AuthorizedUser{
		Username:  entity.Username,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(t.RefreshTokenHours)).Unix(),
	})

	if err != nil {
		return nil, err
	}

	remoteIP, userAgent := getCallerInfo(ctx)
	err = t.SecurityLogService.LogEvent(ctx, entity.UserId, "Login", remoteIP, userAgent)
	if err != nil {
		return nil, err
	}

	t.loginCnt.Inc()

	return &pb.LoginResponse{
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}

func (t *implUIGrpcServer) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if ok {
		t.AuthorizationMiddleware.InvalidateToken(user.Token)
	}

	return &emptypb.Empty{}, nil
}

func (t *implUIGrpcServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (resp *pb.LoginResponse, err error) {
	
	user, err := t.AuthorizationMiddleware.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "Refresh", user.Username)
		}

	}()

	userId, err := t.UserService.GetUserIdByUsername(ctx, user.Username)
	if err == service.ErrUserNotFound {
		err = status.Errorf(codes.NotFound, "username not found")
		return
	}
	if err != nil {
		return
	}

	info, err := t.UserService.GetUser(ctx, userId)
	if err == service.ErrUserNotFound {
		err = status.Errorf(codes.NotFound, "user not found")
		return
	}
	if err != nil {
		return
	}

	roles := make(map[string]bool)
	roles["WEB_USER"] = true
	if info.Role == pb.UserRole_ADMIN {
		roles["WEB_ADMIN"] = true
	}

	token, err := t.AuthorizationMiddleware.GenerateToken(&sprint.AuthorizedUser{
		Username:  user.Username,
		Roles:     roles,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(t.AccessTokenMinutes)).Unix(),
	})

	if err != nil {
		return nil, err
	}

	refreshToken, err := t.AuthorizationMiddleware.GenerateToken(&sprint.AuthorizedUser{
		Username:  user.Username,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(t.RefreshTokenHours)).Unix(),
	})

	if err != nil {
		return
	}

	return &pb.LoginResponse{
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}

func (t *implUIGrpcServer) IsUsernameAvailable(ctx context.Context, req *pb.UsernameRequest) (resp *pb.UsernameResponse, err error) {

	remoteIP, _ := getCallerInfo(ctx)

	var rateLimiter *sprintutils.RateLimiter
	if value, ok := t.usernameLimiterMap.Load(remoteIP); ok {
		if pointer, ok := value.(*sprintutils.RateLimiter); ok {
			rateLimiter = pointer
		}
	}

	if rateLimiter == nil {
		rateLimiter = new(sprintutils.RateLimiter)
		rateLimiter.Limit = time.Second
		t.usernameLimiterMap.Store(remoteIP, rateLimiter)
	}

	resp = new(pb.UsernameResponse)
	resp.Name = req.Name

	err = rateLimiter.Do(func() error {
		resp.Available, resp.NormName, err = t.UserService.IsUsernameAvailable(ctx, req.Name)
		return err
	})

	if err != nil {
		return nil, t.wrapError(err, "UsernameAvailable", remoteIP)
	}

	return
}

func (t *implUIGrpcServer) User(ctx context.Context, _ *emptypb.Empty) (*pb.UserResponse, error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_USER"] {
		return nil, status.Errorf(codes.Unauthenticated, "user not authorized")
	}

	userId, err := t.UserService.GetUserIdByUsername(ctx, user.Username)
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user id not found")
	}
	if err != nil {
		return nil, t.wrapError(err, "User", user.Username)
	}

	info, err := t.UserService.GetUser(ctx, userId)
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, t.wrapError(err, "UserId", userId)
	}

	u := &pb.User{
		UserId:     info.UserId,
		Username:   info.Username,
		FirstName:  info.FirstName,
		MiddleName: info.MiddleName,
		LastName:   info.LastName,
		Email:      info.Email,
		Since:      int64(time.Unix(info.CreTimestamp, 0).Year()),
		Role:       t.getWebUserRole(user),
	}

	return &pb.UserResponse{
		User: u,
	}, nil
}

func (t *implUIGrpcServer) getWebUserRole(user *sprint.AuthorizedUser) string {
	var role string
	if user.Roles["WEB_USER"] {
		role = "USER"
	}
	if user.Roles["WEB_ADMIN"] {
		role = "ADMIN"
	}
	return role
}

func (t *implUIGrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {

	_, err := t.UserService.GetUserIdByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already registered")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "Register", req.Email)
		}

	}()

	if err != service.ErrUserNotFound {
		return nil, err
	}

	entity, err := t.UserService.CreateUser(ctx, req)
	if err == service.ErrUserAlreadyExist {
		return nil, status.Errorf(codes.AlreadyExists, "user already exist")
	}

	subject := fmt.Sprintf("Welcome to %s, %s.", t.WebappName, req.FirstName)
	sender := t.Properties.GetString("mail.sender", "noreply@localhost")

	mail := sprint.Mail{
		Sender:      sender,
		Recipients:   []string{req.Email},
		Subject:      subject,
		TextTemplate: "resources:mail/register_text.tmpl",
		HtmlTemplate: "resources:mail/register_html.tmpl",
		Data:         map[string]interface{} {
			"FirstName": req.FirstName,
			"Project": t.WebappName,
		},
	}

	err = t.MailService.SendMail(&mail, time.Minute, true)
	if err != nil {
		return nil, err
	}

	adminEmail := t.Properties.GetString("webapp.admin", "")
	if adminEmail != "" {

		mail := sprint.Mail{
			Sender:       t.Properties.GetString("mail.sender", "noreply@localhost"),
			Recipients:   []string{adminEmail},
			Subject:      fmt.Sprintf("New user on %s.", t.WebappName),
			TextTemplate: "resources:mail/user_registered_text.tmpl",
			HtmlTemplate: "resources:mail/user_registered_html.tmpl",
			Data:         map[string]interface{} {
				"FirstName": req.FirstName,
				"LastName": req.LastName,
				"Email": req.Email,
				"Project": t.WebappName,
			},
		}

		go t.MailService.SendMail(&mail, time.Minute, false)
	}

	remoteIP, userAgent := getCallerInfo(ctx)
	err = t.SecurityLogService.LogEvent(ctx, entity.UserId, "Registration", remoteIP, userAgent)
	if err != nil {
		return nil, err
	}

	t.registerCnt.Inc()

	return &emptypb.Empty{}, err
}

func (t *implUIGrpcServer) Restore(ctx context.Context, req *pb.RestoreRequest) (*emptypb.Empty, error) {

	resp, err := t.doRestore(ctx, req)
	if err != nil {
		return nil, t.wrapError(err, "Restore", req.Login)
	}

	return resp, nil
}

func (t *implUIGrpcServer) doRestore(ctx context.Context, req *pb.RestoreRequest) (*emptypb.Empty, error) {

	//t.Log.Info("Restore", zap.Any("req", req.String()))

	userId, err := t.UserService.GetUserIdByLogin(ctx, req.Login)
	if err == service.ErrUserNotFound {
		// do nothing, let's make illusion that this email also registered
		t.Log.Info("RestoreUserNotFound", zap.Any("login", req.Login))
		return &emptypb.Empty{}, nil
	}
	if err != nil {
		return nil, err
	}

	entity, err := t.UserService.GetUser(ctx, userId)
	if err == service.ErrUserNotFound {
		// do nothing, let's make illusion that this email also registered
		return &emptypb.Empty{}, nil
	}
	if err != nil {
		return nil, err
	}


	code := strconv.FormatInt(int64(rand.Int31()), 10)
	subject := fmt.Sprintf("%s is %s recover passcode", code, t.WebappName)
	sender := t.Properties.GetString("mail.sender", "noreply@localhost")

	remoteIP, _ := getCallerInfo(ctx)

	mail := sprint.Mail{
		Sender:      sender,
		Recipients:   []string{entity.Email},
		Subject:      subject,
		TextTemplate: "resources:mail/recover_text.tmpl",
		HtmlTemplate: "resources:mail/recover_html.tmpl",
		Data:         map[string]interface{} {
			"Code": code,
			"RemoteIP": remoteIP,
			"Time": time.Now().String(),
			"Project": t.WebappName,
		},
	}

	err = t.MailService.SendMail(&mail, time.Minute, true)
	if err != nil {
		return nil, err
	}

	err = t.UserService.SaveRecoverCode(ctx, req.Login, &pb.RecoverCodeEntity{
		Code:         code,
		RemoteIp:     remoteIP,
		CreTimestamp: time.Now().Unix(),
	}, 60 * 20)

	t.restoreCnt.Inc()

	return &emptypb.Empty{}, err
}

func (t *implUIGrpcServer) Reset(ctx context.Context, req *pb.ResetRequest) (resp *emptypb.Empty, err error) {

	//t.Log.Info("Reset", zap.Any("req", req.String()))

	err = t.UserService.ValidateRecoverCode(ctx, req.Login, req.Code)
	if err == service.ErrInvalidRecoverCode {
		return nil, status.Errorf(codes.InvalidArgument, "wrong recovery code")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "Reset", req.Login)
		}

	}()

	if err != nil {
		return nil, err
	}

	userId, err := t.UserService.GetUserIdByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	email, err := t.UserService.ResetPassword(ctx, userId, req.Password)
	if err != nil {
		return nil, err
	}

	sender := t.Properties.GetString("mail.sender", "noreply@localhost")
	support := t.Properties.GetString("mail.support", "support@localhost")

	subject := fmt.Sprintf("Password reset for %s.", req.Login)
	remoteIP, userAgent := getCallerInfo(ctx)

	err = t.SecurityLogService.LogEvent(ctx, userId, "ResetPassword", remoteIP, userAgent)
	if err != nil {
		return nil, err
	}

	mail := sprint.Mail{
		Sender:      sender,
		Recipients:   []string{email},
		Subject:      subject,
		TextTemplate: "resources:mail/reset_text.tmpl",
		HtmlTemplate: "resources:mail/reset_html.tmpl",
		Data:         map[string]interface{} {
			"RemoteIP": remoteIP,
			"HelpEmail": support,
			"Project": t.WebappName,
		},
	}

	err = t.MailService.SendMail(&mail, time.Minute, true)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *implUIGrpcServer) SecurityLog(ctx context.Context, req *pb.SecurityLogRequest) (resp *pb.SecurityLogResponse, err error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_USER"] {
		return nil, status.Errorf(codes.Unauthenticated, "user not authorized")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "SecurityLog", user.Username)
		}

	}()

	userId, err := t.UserService.GetUserIdByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	var log []*pb.SecurityLogEntity
	err = t.SecurityLogService.EnumEvents(ctx, userId, func(event *pb.SecurityLogEntity) bool {
		log = append(log, event)
		return true
	})
	if err != nil {
		return nil, err
	}

	total := len(log)
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}
	limit := int(req.Limit)

	if offset >= total {
		return &pb.SecurityLogResponse{Total: int32(total)}, nil
	}
	var items []*pb.SecurityLogItem

	for j := total - 1 - offset; j >= 0 && limit > 0; j-- {

		items = append(items, &pb.SecurityLogItem{
			Position: int32(j + 1),
			EventName: log[j].EventName,
			EventTime: log[j].EventTime,
			RemoteIp:  log[j].RemoteIp,
			UserAgent: log[j].UserAgent,
		})

		limit--
	}

	return &pb.SecurityLogResponse{
		Total:   int32(total),
		Items:   items,
	}, nil
}
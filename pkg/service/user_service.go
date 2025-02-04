/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package service

import (
	"context"
	"fmt"
	"go.arpabet.com/store"
	"github.com/pkg/errors"
	"go.arpabet.com/sprint"
	"go.arpabet.com/sprintframework/sprintutils"
	"go.arpabet.com/template/pkg/api"
	"go.arpabet.com/template/pkg/pb"
	"go.arpabet.com/template/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	minUsernameLength = 5
)

type implUserService struct {
	Log                  *zap.Logger                `inject`
	ConfigRepository     sprint.ConfigRepository    `inject`
	HostStore            store.ManagedDataStore     `inject:"bean=host-store"`
	TransactionalManager store.TransactionalManager `inject:"bean=host-store"`

	UserSaltKey   string `value:"user-service.salt-key,default="`
	InitialUserId int    `value:"user-service.initial-id,default=27483984961"` // u00001
}

func UserService() api.UserService {
	return &implUserService{}
}

func (t *implUserService) PostConstruct() (err error) {
	if t.UserSaltKey == "" {
		t.UserSaltKey, err = sprintutils.GenerateToken()
		if err != nil {
			return errors.Errorf("generate token error, %v", err)
		}
		err = t.ConfigRepository.Set("user-service.salt-key", t.UserSaltKey)
		return err
	}
	return nil
}

func (t *implUserService) CreateUser(ctx context.Context, req *pb.RegisterRequest) (user *pb.UserEntity, err error) {

	req.Username = utils.NormalizeUsername(req.Username)
	if req.Username == "" {
		return nil, errors.New("username is empty")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	if req.Email == "" {
		return nil, errors.New("user email is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	user = new(pb.UserEntity)
	err = t.HostStore.Get(ctx).ByKey("%s:user", req.Email).ToProto(user)
	if err != nil {
		return nil, err
	}
	if user.Email != "" {
		return nil, ErrUserAlreadyExist
	}

	if req.Password == "" {
		return nil, errors.New("user password is empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(t.UserSaltKey+req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := pb.UserRole_USER
	if has, err := t.hasUsers(ctx); err != nil {
		return nil, err
	} else if !has {
		role = pb.UserRole_ADMIN
	}

	userId, err := t.GenerateUserId(ctx)
	if err != nil {
		return nil, err
	}

	user = &pb.UserEntity{
		UserId:       userId,
		Username:     req.Username,
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreTimestamp: time.Now().Unix(),
		Role:         role,
	}

	err = t.HostStore.Set(ctx).ByKey("%s:user", userId).Proto(user)
	if err != nil {
		return nil, err
	}

	// back reference
	err = t.HostStore.Set(ctx).ByKey("user:%s", userId).String(userId)
	if err != nil {
		return nil, err
	}

	// username index
	err = t.HostStore.Set(ctx).ByKey("username:%s", req.Username).String(userId)

	// email index
	err = t.HostStore.Set(ctx).ByKey("email:%s", req.Email).String(userId)

	return user, err
}

func (t *implUserService) hasUsers(ctx context.Context) (has bool, err error) {
	err = t.EnumUsers(ctx, func(user *pb.UserEntity) bool {
		has = true
		return false
	})
	return
}

func (t *implUserService) GenerateUserId(ctx context.Context) (string, error) {
	for {
		num, err := t.HostStore.Increment(ctx).ByKey("user-next-id").WithInitialValue(int64(t.InitialUserId)).WithDelta(1).Do()
		if err != nil {
			return "", err
		}
		id := sprintutils.EncodeId(uint64(num))
		value, err := t.HostStore.Get(ctx).ByKey("user:%s", id).ToString()
		if err != nil {
			return "", err
		}
		if value == "" {
			return id, nil
		}
	}
}

func (t *implUserService) ResetPassword(ctx context.Context, userId string, newPassword string) (email string, err error) {

	if newPassword == "" {
		return "", errors.New("new password is empty")
	}

	err = t.DoWithUser(ctx, userId, func(user *pb.UserEntity) error {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(t.UserSaltKey+newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user.PasswordHash = hashedPassword

		email = user.Email
		return nil
	})

	return
}

func (t *implUserService) AuthenticateUser(ctx context.Context, login, password string) (*pb.UserEntity, error) {

	login = utils.NormalizeLogin(login)
	if login == "" {
		return nil, errors.New("empty login")
	}

	userId, err := t.HostStore.Get(ctx).ByKey("username:%s", login).ToString()
	if err != nil {
		return nil, err
	}
	if userId == "" {
		userId, err = t.HostStore.Get(ctx).ByKey("email:%s", login).ToString()
		if err != nil {
			return nil, err
		}
		if userId == "" {
			userId = login
		}
	}

	user := new(pb.UserEntity)
	err = t.HostStore.Get(ctx).ByKey("%s:user", userId).ToProto(user)
	if err != nil {
		return nil, err
	}
	if user.UserId != userId {
		return nil, ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(t.UserSaltKey+password))
	if err != nil {
		return user, ErrUserInvalidPassword
	}
	return user, nil
}

func (t *implUserService) IsUsernameAvailable(ctx context.Context, username string) (bool, string, error) {

	normName := utils.NormalizeUsername(username)
	if len(normName) < minUsernameLength {
		return false, normName, nil
	}
	userId, err := t.HostStore.Get(ctx).ByKey("username:%s", normName).ToString()
	return userId == "", normName, err
}

func (t *implUserService) GetUser(ctx context.Context, userId string) (*pb.UserEntity, error) {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return nil, errors.New("user id is empty")
	}

	user := new(pb.UserEntity)
	err := t.HostStore.Get(ctx).ByKey("%s:user", userId).ToProto(user)
	if err != nil {
		return nil, err
	}
	if user.UserId == "" {
		return user, ErrUserNotFound
	}
	if user.UserId != userId {
		t.Log.Error("FindUser",
			zap.String("Value", user.String()),
			zap.String("User", userId),
			zap.Error(ErrIntegrityDB))
		return nil, ErrIntegrityDB
	}
	return user, nil
}

func (t *implUserService) GetUserIdByLogin(ctx context.Context, login string) (string, error) {

	login = utils.NormalizeLogin(login)
	if login == "" {
		return "", errors.New("empty login")
	}

	userId, err := t.HostStore.Get(ctx).ByKey("username:%s", login).ToString()
	if err != nil {
		return "", err
	}
	if userId != "" {
		return userId, nil
	}

	userId, err = t.HostStore.Get(ctx).ByKey("email:%s", login).ToString()
	if err != nil {
		return "", err
	}
	if userId != "" {
		return userId, nil
	}

	return "", ErrUserNotFound
}

func (t *implUserService) GetUserIdByEmail(ctx context.Context, email string) (string, error) {

	email = utils.NormalizeEmail(email)
	if email == "" {
		return "", errors.New("empty user email")
	}

	userId, err := t.HostStore.Get(ctx).ByKey("email:%s", email).ToString()
	if err != nil {
		return "", err
	}
	if userId == "" {
		return "", ErrUserNotFound
	}

	return userId, nil
}

func (t *implUserService) GetUserIdByUsername(ctx context.Context, username string) (string, error) {

	username = utils.NormalizeUsername(username)
	if username == "" {
		return "", errors.New("empty username")
	}

	userId, err := t.HostStore.Get(ctx).ByKey("username:%s", username).ToString()
	if err != nil {
		return "", err
	}
	if userId == "" {
		return "", ErrUserNotFound
	}

	return userId, nil
}

func (t *implUserService) SaveUser(ctx context.Context, user *pb.UserEntity) (err error) {

	user.UserId = utils.NormalizeUserId(user.UserId)
	if user.UserId == "" {
		return errors.New("user id is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	oldUser := new(pb.UserEntity)
	err = t.HostStore.Get(ctx).ByKey("%s:user", user.UserId).ToProto(oldUser)
	if err != nil {
		return err
	}

	// if found
	if oldUser.UserId == user.UserId {
		if oldUser.Username != user.Username {
			err = t.HostStore.Remove(ctx).ByKey("username:%s", oldUser.Username).Do()
			if err != nil {
				return err
			}
		}
		if oldUser.Email != user.Email {
			err = t.HostStore.Remove(ctx).ByKey("email:%s", oldUser.Email).Do()
			if err != nil {
				return err
			}
		}
	}

	err = t.HostStore.Set(ctx).ByKey("%s:user", user.UserId).Proto(user)
	if err != nil {
		return err
	}

	// username index check
	usedUserId, err := t.HostStore.Get(ctx).ByKey("username:%s", user.Username).ToString()
	if err != nil {
		return err
	}
	if usedUserId != user.UserId {
		return errors.Errorf("username '%s' is already used by user '%s'", user.Username, usedUserId)
	}

	// username index
	err = t.HostStore.Set(ctx).ByKey("username:%s", user.Username).String(user.UserId)
	if err != nil {
		return err
	}

	// email index check
	usedUserId, err = t.HostStore.Get(ctx).ByKey("email:%s", user.Email).ToString()
	if err != nil {
		return err
	}
	if usedUserId != user.UserId {
		return errors.Errorf("email '%s' is already used by user '%s'", user.Email, usedUserId)
	}

	// email index
	err = t.HostStore.Set(ctx).ByKey("email:%s", user.Email).String(user.UserId)
	if err != nil {
		return err
	}

	// back reference
	err = t.HostStore.Set(ctx).ByKey("user:%s", user.UserId).String(user.UserId)
	return
}

func (t *implUserService) RemoveUser(ctx context.Context, userId string) (err error) {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return errors.New("user id is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	user, err := t.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	err = t.doRemoveUser(ctx, user)
	if err != nil {
		return err
	}

	return
}

func (t *implUserService) doRemoveUser(ctx context.Context, user *pb.UserEntity) (err error) {

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	// remove user object
	err = t.HostStore.Remove(ctx).ByKey("%s:user", user.UserId).Do()
	if err != nil {
		return err
	}

	// remove back references
	err = t.HostStore.Remove(ctx).ByKey("user:%s", user.UserId).Do()
	if err != nil {
		return err
	}

	err = t.HostStore.Remove(ctx).ByKey("username:%s", user.Username).Do()
	if err != nil {
		return err
	}

	err = t.HostStore.Remove(ctx).ByKey("email:%s", user.Email).Do()
	if err != nil {
		return err
	}

	return nil
}

func (t *implUserService) DropUserContent(ctx context.Context, userId string) error {
	prefix := fmt.Sprintf("%s:", userId)
	return t.HostStore.DropWithPrefix([]byte(prefix))
}

func (t *implUserService) EnumUsers(ctx context.Context, cb func(user *pb.UserEntity) bool) error {

	return t.HostStore.Enumerate(ctx).
		ByPrefix("user:").
		WithBatchSize(BatchSize).
		Do(func(entry *store.RawEntry) bool {
			userId := string(entry.Value)
			user := new(pb.UserEntity)
			err := t.HostStore.Get(context.Background()).ByKey("%s:user", userId).ToProto(user)
			if err == nil && user.Email != "" {
				return cb(user)
			} else if user.Email == "" {
				t.Log.Warn("UserNotFound", zap.String("backwardKey", string(entry.Key)), zap.String("userId", userId))
			} else {
				t.Log.Warn("EnumUsers", zap.Error(err), zap.String("backwardKey", string(entry.Key)), zap.String("userId", userId))
			}
			return true
		})

}

func (t *implUserService) DoWithUser(ctx context.Context, userId string, cb func(user *pb.UserEntity) error) (err error) {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return errors.New("user id is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	user, err := t.GetUser(ctx, userId)
	if err == ErrUserNotFound {
		return err
	}
	if err != nil {
		return errors.Errorf("load user '%s', %v", userId, err)
	}
	savedUsername := user.Username
	savedEmail := user.Email

	err = cb(user)
	if err != nil {
		return err
	}

	if savedUsername != user.Username {

		usedUserId, err := t.HostStore.Get(ctx).ByKey("username:%s", user.Username).ToString()
		if err != nil {
			return err
		}
		if usedUserId != "" {
			return errors.Errorf("username '%s' is already used by user '%s'", user.Username, usedUserId)
		}

		err = t.HostStore.Set(ctx).ByKey("username:%s", user.Username).String(userId)
		if err != nil {
			return err
		}

		err = t.HostStore.Remove(ctx).ByKey("username:%s", savedUsername).Do()
		if err != nil {
			return err
		}
	}

	if savedEmail != user.Email {

		usedUserId, err := t.HostStore.Get(ctx).ByKey("email:%s", user.Email).ToString()
		if err != nil {
			return err
		}
		if usedUserId != "" {
			return errors.Errorf("email '%s' is already used by user '%s'", user.Email, usedUserId)
		}

		err = t.HostStore.Set(ctx).ByKey("email:%s", user.Email).String(userId)
		if err != nil {
			return err
		}

		err = t.HostStore.Remove(ctx).ByKey("email:%s", savedEmail).Do()
		if err != nil {
			return err
		}
	}

	err = t.HostStore.Set(ctx).ByKey("%s:user", user.UserId).Proto(user)
	if err != nil {
		return err
	}

	return nil
}

func (t *implUserService) DumpUser(ctx context.Context, userId string, cb func(entry *store.RawEntry) bool) error {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return errors.New("user id is empty")
	}

	return t.HostStore.Enumerate(ctx).
		ByPrefix("%s:", userId).
		WithBatchSize(BatchSize).
		Do(func(entry *store.RawEntry) bool {
			return cb(entry)
		})

}

func (t *implUserService) SaveRecoverCode(ctx context.Context, login string, rc *pb.RecoverCodeEntity, ttlSeconds int) error {

	login = utils.NormalizeLogin(login)
	if login == "" {
		return errors.New("login is empty")
	}

	return t.HostStore.Set(ctx).ByKey("recover:login:%s", login).WithTtl(ttlSeconds).Proto(rc)
}

func (t *implUserService) ValidateRecoverCode(ctx context.Context, login string, code string) error {

	login = utils.NormalizeLogin(login)
	if login == "" {
		return errors.New("login is empty")
	}

	code = utils.NormalizeCode(code)
	if code == "" {
		return errors.New("user code is empty")
	}

	rc := new(pb.RecoverCodeEntity)
	err := t.HostStore.Get(ctx).ByKey("recover:login:%s", login).ToProto(rc)
	if err != nil {
		return err
	}

	if rc.Code != code {
		return ErrInvalidRecoverCode
	}

	return nil

}

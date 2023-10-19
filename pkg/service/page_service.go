/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package service

import (
	"context"
	"github.com/keyvalstore/store"
	"github.com/pkg/errors"
	"github.com/sprintframework/template/pkg/api"
	"github.com/sprintframework/template/pkg/pb"
	"github.com/sprintframework/template/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

type implPageService struct {
	Log                  *zap.Logger                `inject`
	HostStore            store.DataStore            `inject:"bean=host-store"`
	TransactionalManager store.TransactionalManager `inject:"bean=host-store"`
}

func PageService() api.PageService {
	return &implPageService{}
}

func (t *implPageService) GetPage(ctx context.Context, name string) (*pb.PageEntity, error) {

	name = utils.NormalizePageId(name)
	if name == "" {
		return nil, errors.New("page name is empty")
	}

	page := new(pb.PageEntity)
	err := t.HostStore.Get(ctx).ByKey("page:%s", name).ToProto(page)
	if err != nil {
		return nil, err
	}
	if page.Name == "" {
		return nil, ErrPageNotFound
	}
	if page.Name != name {
		t.Log.Error("FindPage",
			zap.String("value", page.String()),
			zap.String("name", name),
			zap.Error(ErrIntegrityDB))
		return nil, ErrIntegrityDB
	}
	return page, nil

}

func (t *implPageService) CreatePage(ctx context.Context, newPage *pb.AdminPage) (err error) {

	newPage.Name = utils.NormalizePageId(newPage.Name)
	if newPage.Name == "" {
		return errors.New("new page name is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	entity := new(pb.PageEntity)
	err = t.HostStore.Get(ctx).ByKey("page:%s", newPage.Name).ToProto(entity)
	if err != nil {
		return
	}

	if entity.Name != "" {
		err = errors.Errorf("nowrap: page '%s' already exist", newPage.Name)
		return
	}

	contentType, err := t.parseContentType(newPage.ContentType)
	if err != nil {
		err = errors.Errorf("nowrap: invalid content type '%s'", newPage.ContentType)
	}

	entity = &pb.PageEntity{
		Name:         newPage.Name,
		Title:        newPage.Title,
		Content:      newPage.Content,
		ContentType:  contentType,
		CreTimestamp: time.Now().Unix(),
	}

	err = t.HostStore.Set(ctx).ByKey("page:%s", newPage.Name).Proto(entity)
	return

}

func (t *implPageService) UpdatePage(ctx context.Context, updatingPage *pb.AdminPage) (err error) {

	updatingPage.Name = utils.NormalizePageId(updatingPage.Name)
	if updatingPage.Name == "" {
		return errors.New("updating page name is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	if updatingPage.Name != updatingPage.Prev && updatingPage.Prev != "" {

		entity := new(pb.PageEntity)
		err = t.HostStore.Get(ctx).ByKey("page:%s", updatingPage.Name).ToProto(entity)
		if err != nil {
			return
		}

		if entity.Name != "" {
			err = errors.Errorf("nowrap: page '%s' already exist", updatingPage.Name)
			return
		}

		err = t.HostStore.Remove(ctx).ByKey("page:%s", updatingPage.Prev).Do()
		if err != nil {
			return
		}

	}

	contentType, err := t.parseContentType(updatingPage.ContentType)
	if err != nil {
		err = errors.Errorf("nowrap: invalid content type '%s'", updatingPage.ContentType)
	}

	entity := &pb.PageEntity{
		Name:         updatingPage.Name,
		Title:        updatingPage.Title,
		Content:      updatingPage.Content,
		ContentType:  contentType,
		CreTimestamp: time.Now().Unix(),
	}

	err = t.HostStore.Set(ctx).ByKey("page:%s", updatingPage.Name).Proto(entity)
	return

}

func (t *implPageService) RemovePage(ctx context.Context, name string) error {

	name = utils.NormalizePageId(name)
	if name == "" {
		return errors.New("page name is empty")
	}

	return t.HostStore.Remove(ctx).ByKey("page:%s", name).Do()
}

func (t *implPageService) EnumPages(ctx context.Context, cb func(page *pb.PageEntity) bool) error {

	return t.HostStore.Enumerate(ctx).
		ByPrefix("page:").
		WithBatchSize(BatchSize).
		DoProto(func() proto.Message {
			return new(pb.PageEntity)
		}, func(entry *store.ProtoEntry) bool {
			if v, ok := entry.Value.(*pb.PageEntity); ok {
				return cb(v)
			}
			return true
		})

}

func (t *implPageService) parseContentType(ct string) (pb.ContentType, error) {
	contentType := pb.ContentType_MARKDOWN
	switch strings.ToUpper(strings.TrimSpace(ct)) {
	case "MARKDOWN":
		contentType = pb.ContentType_MARKDOWN
	case "HTML":
		contentType = pb.ContentType_HTML
	default:
		return 0, errors.Errorf("invalid content type '%s'", ct)
	}
	return contentType, nil
}

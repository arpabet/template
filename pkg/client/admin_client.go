/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */


package client

import (
	"context"
	"github.com/sprintframework/template/pkg/api"
	"github.com/sprintframework/template/pkg/pb"
	"google.golang.org/grpc"
	"sync"
)

type implAdminClient struct {
	GrpcConn   *grpc.ClientConn                `inject`
	client     pb.AdminServiceClient
	closeOnce  sync.Once
}

func AdminClient() api.AdminClient {
	return &implAdminClient{}
}

func (t *implAdminClient) PostConstruct() error {
	t.client = pb.NewAdminServiceClient(t.GrpcConn)
	return nil
}

func (t *implAdminClient) AdminCommand(command string, args []string) (string, error) {

	req := &pb.Command {
		Command: command,
		Args: args,
	}

	if resp, err := t.client.AdminRun(context.Background(), req); err != nil {
		return "", err
	} else {
		return resp.Content, nil
	}
}

func (t *implAdminClient) Destroy() (err error) {
	t.closeOnce.Do(func() {
		if t.GrpcConn != nil {
			err = t.GrpcConn.Close()
		}
	})
	return
}



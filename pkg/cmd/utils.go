/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package cmd

import (
	"github.com/codeallergy/glue"
	"github.com/pkg/errors"
	"github.com/sprintframework/sprint"
	"github.com/sprintframework/template/pkg/api"
	"reflect"
)

func doWithAdminClient(parent glue.Context, cb func(client api.AdminClient) error) error {

	return sprint.DoWithClient(parent, sprint.ControlClientRole, api.AdminClientClass, func(instance interface{}) error {

		if client, ok := instance.(api.AdminClient); ok {
			return cb(client)
		} else {
			return errors.Errorf("invalid object '%v' found instead of api.AdminClient in client context: ", reflect.TypeOf(instance))
		}

	})
}
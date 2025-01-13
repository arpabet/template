/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package cmd

import (
	"go.arpabet.com/glue"
	"github.com/pkg/errors"
	"go.arpabet.com/sprint"
	"go.arpabet.com/template/pkg/api"
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
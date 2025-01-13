/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */


package api

import (
	"go.arpabet.com/glue"
	"reflect"
)

var AdminClientClass = reflect.TypeOf((*AdminClient)(nil)).Elem()

type AdminClient interface {
	glue.InitializingBean
	glue.DisposableBean

	AdminCommand(command string, args []string) (string, error)

}

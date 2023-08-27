/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */


package api

import (
	"github.com/codeallergy/glue"
	"reflect"
)

var AdminClientClass = reflect.TypeOf((*AdminClient)(nil)).Elem()

type AdminClient interface {
	glue.InitializingBean
	glue.DisposableBean

	AdminCommand(command string, args []string) (string, error)

}

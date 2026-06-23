/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package api

import (
	"go.arpabet.com/glue"
	"go.arpabet.com/sprint"
	"reflect"
)

var GRPCServerClass = reflect.TypeOf((*GRPCServer)(nil)).Elem()

type GRPCServer interface {
	glue.InitializingBean
	sprint.Component
}



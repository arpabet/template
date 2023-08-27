/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package cmd

import (
	"github.com/codeallergy/glue"
	"github.com/sprintframework/template/pkg/api"
	"github.com/sprintframework/sprint"
	"github.com/pkg/errors"
	"log"
)

func doWithAdminClient(parent glue.Context, cb func(client api.AdminClient) error) error {

	var verbose bool
	list := parent.Bean(sprint.ApplicationFlagsClass, glue.DefaultLevel)
	if len(list) > 0 {
		flags := list[0].Object().(sprint.ApplicationFlags)
		if flags.Verbose() {
			verbose = true
	}
	}

	list = parent.Bean(sprint.ClientScannerClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("application context should have one client scanner, but found '%d'", len(list))
	}
	bean := list[0]

	scanner, ok := bean.Object().(sprint.ClientScanner)
	if !ok {
		return errors.Errorf("invalid object '%v' found instead of sprint.ClientScanner in application context", bean.Class())
	}

	beans := scanner.ClientBeans()
	if verbose {
		verbose := glue.Verbose{ Log: log.Default() }
		beans = append([]interface{}{verbose}, beans...)
	}

	ctx, err := parent.Extend(beans...)
	if err != nil {
		return err
	}
	defer ctx.Close()

	list = ctx.Bean(api.AdminClientClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("client context shoulw have one api.AdminClient inference, but found '%d'", len(list))
	}
	bean = list[0]

	if client, ok := bean.Object().(api.AdminClient); ok {
		return cb(client)
	} else {
		return errors.Errorf("invalid object '%v' found instead of api.AdminClient in client context", bean.Class())
	}

}
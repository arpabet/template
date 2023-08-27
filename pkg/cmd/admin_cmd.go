/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package cmd

import (
	"github.com/pkg/errors"
	"github.com/codeallergy/glue"
	"github.com/sprintframework/template/pkg/api"
	"github.com/sprintframework/sprint"
)

type implAdminCommand struct {
	Context           glue.Context             `inject`
	Application       sprint.Application        `inject`
	ApplicationFlags   sprint.ApplicationFlags   `inject`
}

func AdminCommand() sprint.Command {
	return &implAdminCommand{}
}

func (t *implAdminCommand) BeanName() string {
	return "admin"
}

func (t *implAdminCommand) Desc() string {
	return "admin commands: [list, add, remove]"
}

func (t *implAdminCommand) Run(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("admin command needs argument, %s", t.Desc())
	}
	cmd := args[0]
	args = args[1:]

	return doWithAdminClient(t.Context, func(client api.AdminClient) error {
		content, err := client.AdminCommand(cmd, args)
		if err == nil {
			println(content)
		}
		return err
	})

}

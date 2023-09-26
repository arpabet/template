//go:generate make generate
//go:generate python3 gtag.py MYGTAG assets/

/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/codeallergy/glue"
	"github.com/sprintframework/certmod"
	"github.com/sprintframework/dnsmod"
	"github.com/sprintframework/natmod"
	"github.com/sprintframework/sealmod"
	"github.com/sprintframework/sprint"
	"github.com/sprintframework/template/pkg/assets"
	"github.com/sprintframework/template/pkg/assetsgz"
	"github.com/sprintframework/template/pkg/client"
	"github.com/sprintframework/template/pkg/cmd"
	"github.com/sprintframework/template/pkg/resources"
	"github.com/sprintframework/template/pkg/server"
	"github.com/sprintframework/template/pkg/service"
	"github.com/sprintframework/sprintframework/sprintapp"
	"github.com/sprintframework/sprintframework/sprintclient"
	"github.com/sprintframework/sprintframework/sprintcmd"
	"github.com/sprintframework/sprintframework/sprintcore"
	"github.com/sprintframework/sprintframework/sprintserver"
	"os"
	"time"
)

var (
	Version string
	Build   string
)

var AppAssets = &glue.ResourceSource{
	Name: "assets",
	AssetNames: assets.AssetNames(),
	AssetFiles: assets.AssetFile(),
}

var AppGzipAssets = &glue.ResourceSource{
	Name: "assets-gzip",
	AssetNames: assetsgz.AssetNames(),
	AssetFiles: assetsgz.AssetFile(),
}

var AppResources = &glue.ResourceSource{
	Name: "resources",
	AssetNames: resources.AssetNames(),
	AssetFiles: resources.AssetFile(),
}

func doMain() (err error) {

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = errors.Errorf("%v", v)
			}
		}
	}()

	beans := []interface{} {
		sprintapp.ApplicationBeans,
		sprintcmd.ApplicationCommands,
		cmd.Commands,

		AppAssets,
		AppGzipAssets,
		AppResources,

		glue.Child(sprint.CoreRole,
			sprintcore.CoreServices,
			natmod.Scanner(),
			dnsmod.Scanner(),
			sealmod.Scanner(),
			certmod.Scanner(),
			sprintcore.BadgerStoreFactory("config-store"),
			sprintcore.BadgerStoreFactory("host-store"),
			sprintcore.LumberjackFactory(),
			sprintcore.AutoupdateService(),
			service.UserService(),
			service.SecurityLogService(),
			service.PageService(),

			glue.Child(sprint.ServerRole,
				sprintserver.GrpcServerScanner("control-grpc-server"),
				sprintserver.ControlServer(),
				server.UIGrpcServer(),
				sprintserver.HttpServerFactory("control-gateway-server"),
				sprintserver.TlsConfigFactory("tls-config"),
			),

		),
		glue.Child(sprint.ControlClientRole,
			sprintclient.ControlClientBeans,
			sprintclient.AnyTlsConfigFactory("tls-config"),
			client.AdminClient(),
		),
	}

	return sprintapp.Application("template",
		sprintapp.WithVersion(Version),
		sprintapp.WithBuild(Build),
		sprintapp.WithBeans(beans)).
		Run(os.Args[1:])

}

func main() {

	if err := doMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(100 * time.Millisecond)
}

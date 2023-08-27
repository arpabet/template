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
	"github.com/sprintframework/template/pkg/assets"
	"github.com/sprintframework/template/pkg/assetsgz"
	"github.com/sprintframework/template/pkg/client"
	"github.com/sprintframework/template/pkg/cmd"
	"github.com/sprintframework/template/pkg/resources"
	"github.com/sprintframework/template/pkg/server"
	"github.com/sprintframework/template/pkg/service"
	"github.com/sprintframework/sprintframework/pkg/app"
	sprintclient "github.com/sprintframework/sprintframework/pkg/client"
	sprintcmd "github.com/sprintframework/sprintframework/pkg/cmd"
	sprintcore "github.com/sprintframework/sprintframework/pkg/core"
	sprintserver "github.com/sprintframework/sprintframework/pkg/server"
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

	return app.Application("template",
		app.WithVersion(Version),
		app.WithBuild(Build),
		app.Beans(app.DefaultApplicationBeans, sprintcmd.DefaultCommands, AppAssets, AppGzipAssets, AppResources, cmd.Commands),
		app.Core(sprintcore.CoreScanner(
			sprintcore.BadgerStorageFactory("config-storage"),
			sprintcore.BadgerStorageFactory("host-storage"),
			sprintcore.LumberjackFactory(),
			sprintcore.AutoupdateService(),
			service.UserService(),
			service.SecurityLogService(),
			service.PageService(),
		)),
		app.Server(sprintserver.ServerScanner(
			sprintserver.AuthorizationMiddleware(),
			sprintserver.GrpcServerFactory("control-grpc-server"),
			sprintserver.ControlServer(),
			server.UIGrpcServer(),
			sprintserver.HttpServerFactory("control-gateway-server"),
			sprintserver.TlsConfigFactory("tls-config"),
		)),
		app.Client(sprintclient.ControlClientScanner(
			sprintclient.AnyTlsConfigFactory("tls-config"),
			client.AdminClient(),
		)),
	).Run(os.Args[1:])

}

func main() {

	if err := doMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(100 * time.Millisecond)
}

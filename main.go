//go:generate npm run generate --prefix webapp
//go:generate python3 gtag.py MYGTAG assets/

/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package main

import (
	"fmt"
	"go.arpabet.com/glue"
	"go.arpabet.com/sprint/certmod"
	"go.arpabet.com/sprint/dnsmod"
	"go.arpabet.com/sprint/natmod"
	"go.arpabet.com/sprint/sealmod"
	"go.arpabet.com/sprint/sprint"
	"go.arpabet.com/sprint/sprintframework/sprintapp"
	"go.arpabet.com/sprint/sprintframework/sprintclient"
	"go.arpabet.com/sprint/sprintframework/sprintcmd"
	"go.arpabet.com/sprint/sprintframework/sprintcore"
	"go.arpabet.com/sprint/sprintframework/sprintserver"
	"go.arpabet.com/sprint/sprintframework/sprintutils"
	"go.arpabet.com/template/pkg/assets"
	"go.arpabet.com/template/pkg/assetsgz"
	"go.arpabet.com/template/pkg/client"
	"go.arpabet.com/template/pkg/cmd"
	"go.arpabet.com/template/pkg/resources"
	"go.arpabet.com/template/pkg/server"
	"go.arpabet.com/template/pkg/service"
	"os"
	"time"
)

var (
	Version string
	Build   string
)

var AppAssets = &glue.ResourceSource{
	Name:       "assets",
	AssetNames: assets.AssetNames(),
	AssetFiles: assets.AssetFile(),
}

var AppGzipAssets = &glue.ResourceSource{
	Name:       "assets-gzip",
	AssetNames: assetsgz.AssetNames(),
	AssetFiles: assetsgz.AssetFile(),
}

var AppResources = &glue.ResourceSource{
	Name:       "resources",
	AssetNames: resources.AssetNames(),
	AssetFiles: resources.AssetFile(),
}

func doMain() (err error) {

	sprintutils.PanicToError(&err)

	beans := []interface{}{
		sprintapp.ApplicationBeans,
		sprintcmd.ApplicationCommands,
		cmd.Commands,

		AppAssets,
		AppGzipAssets,
		AppResources,

		glue.Child(sprint.CoreRole,
			sprintcore.CoreServices,
			natmod.NatServices,
			dnsmod.DNSServices,
			sealmod.SealServices,
			certmod.CertServices,
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

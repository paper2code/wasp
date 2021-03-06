package main

import (
	"github.com/iotaledger/hive.go/node"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/plugins/banner"
	"github.com/iotaledger/wasp/plugins/cli"
	"github.com/iotaledger/wasp/plugins/committees"
	"github.com/iotaledger/wasp/plugins/config"
	"github.com/iotaledger/wasp/plugins/dashboard"
	"github.com/iotaledger/wasp/plugins/database"
	"github.com/iotaledger/wasp/plugins/dispatcher"
	"github.com/iotaledger/wasp/plugins/gracefulshutdown"
	"github.com/iotaledger/wasp/plugins/logger"
	"github.com/iotaledger/wasp/plugins/nodeconn"
	"github.com/iotaledger/wasp/plugins/peering"
	"github.com/iotaledger/wasp/plugins/publisher"
	"github.com/iotaledger/wasp/plugins/runvm"
	"github.com/iotaledger/wasp/plugins/testplugins/nodeping"
	"github.com/iotaledger/wasp/plugins/testplugins/roundtrip"
	"github.com/iotaledger/wasp/plugins/webapi"
)

func main() {
	registry.InitFlags()
	parameters.InitFlags()

	plugins := node.Plugins(
		banner.Init(),
		config.Init(),
		logger.Init(),
		gracefulshutdown.Init(),
		webapi.Init(),
		cli.Init(),
		database.Init(),
		peering.Init(),
		nodeconn.Init(),
		dispatcher.Init(),
		committees.Init(),
		runvm.Init(),
		publisher.Init(),
		dashboard.Init(),
	)

	testPlugins := node.Plugins(
		roundtrip.Init(),
		nodeping.Init(),
	)

	node.Run(
		plugins,
		testPlugins,
	)
}

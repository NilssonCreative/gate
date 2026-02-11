package main

import (
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/plugins/pluginmsg"
)

func main() {

	proxy.Plugins = append(proxy.Plugins,
		pluginmsg.Plugin,
	)

	gate.Execute()
}
